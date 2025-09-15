"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import "./main.css";
import { fetchFiles, FileInfo, sendToTrash } from "@/lib/api";
import { Card } from "@/components/Card";
import { ActionsBar } from "@/components/ActionsBar";

type SwipeDir = "left" | "right" | null;

export default function HomePage() {
  const [files, setFiles] = useState<FileInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Swipen tilat
  const [idx, setIdx] = useState(0);               // mikä kortti on päällimmäisenä
  const [dragX, setDragX] = useState(0);           // raahausetäisyys
  const [dragging, setDragging] = useState(false);
  const [anim, setAnim] = useState<SwipeDir>(null);
  const startXRef = useRef<number | null>(null);
  const threshold = 120;                            // kuinka pitkälle pitää vetää ennen hyväksyntää

  // Hae tiedostot
  useEffect(() => {
    fetchFiles("~/Downloads")
      .then((list) => setFiles(list))
      .catch((e) => setError(String(e)))
      .finally(() => setLoading(false));
  }, []);

  const current = files[idx] || null;
  const remaining = Math.max(0, files.length - idx - 1);

  // Onko tiedosto kuva?
  const isImage = useCallback((f: FileInfo | null) => {
    if (!f) return false;
    return /\.(png|jpe?g|gif|webp|bmp|svg)$/i.test(f.ext || f.name);
  }, []);

  // Jos sinulla on esikatselu-endpoint, määritä tähän
  const getPreviewUrl = useCallback((f: FileInfo) => {
    // Esimerkki: `${process.env.NEXT_PUBLIC_BACKEND_URL}/api/open?path=${encodeURIComponent(f.path)}`
    return null; // palauta null jos ei ole esikatselua tarjolla
  }, []);

  // Roskikseen (vasen)
  const trashCurrent = useCallback(async () => {
    const f = files[idx];
    if (!f) return;
    try {
      await sendToTrash(f.path);
      // Poistetaan listasta ja pidetään sama idx (seuraava siirtyy tilalle)
      setFiles((prev) => prev.filter((x, i) => i !== idx));
      // Älä kasvattele idx:ää jos poistit nykyisen; korttipino kompressoituu
      setAnim(null);
      setDragX(0);
    } catch (e) {
      alert("Roskikseen siirto epäonnistui: " + e);
    }
  }, [files, idx]);

  // Säilytä/ohita (oikea)
  const keepCurrent = useCallback(() => {
    if (idx < files.length) {
      setIdx((i) => Math.min(i + 1, files.length)); // seuraava kortti
      setAnim(null);
      setDragX(0);
    }
  }, [files.length, idx]);

  // Pyyhkäisyn käsittely
  const onPointerDown = (e: React.PointerEvent) => {
    setDragging(true);
    startXRef.current = e.clientX;
    (e.target as HTMLElement).setPointerCapture(e.pointerId);
  };
  const onPointerMove = (e: React.PointerEvent) => {
    if (!dragging || startXRef.current == null) return;
    const dx = e.clientX - startXRef.current;
    setDragX(dx);
  };
  const onPointerUp = async () => {
    setDragging(false);
    if (Math.abs(dragX) > threshold) {
      const dir: SwipeDir = dragX < 0 ? "left" : "right";
      setAnim(dir);
      // Pieni viive, että animaatio ehtii
      if (dir === "left") {
        await new Promise((r) => setTimeout(r, 140));
        await trashCurrent();
      } else {
        await new Promise((r) => setTimeout(r, 140));
        keepCurrent();
      }
    } else {
      // palaa keskelle
      setAnim(null);
      setDragX(0);
    }
  };

  // Nuolet & pikanäppäimet
  useEffect(() => {
    const onKey = async (e: KeyboardEvent) => {
      if (!current) return;
      if (e.key === "ArrowLeft") {
        setAnim("left");
        await new Promise((r) => setTimeout(r, 120));
        await trashCurrent();
      } else if (e.key === "ArrowRight") {
        setAnim("right");
        await new Promise((r) => setTimeout(r, 120));
        keepCurrent();
      }
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [current, trashCurrent, keepCurrent]);

  if (loading) return <p className="center muted">Ladataan…</p>;
  if (error) return <p className="center error">Virhe: {error}</p>;
  if (!current) {
    return (
      <main className="stage">
        <h1 className="title">Tiedostot</h1>
        <p className="center">Ei enempää kortteja. Kansiot nuoltu puhtaaksi. ✅</p>
      </main>
    );
  }

  // Kortin inline-tyyli raahaukseen/rotaatioon
  const rotation = Math.max(-12, Math.min(12, dragX / 20));
  const cardStyle: React.CSSProperties =
    dragging || dragX !== 0 || anim
      ? {
          transform:
            anim === "left"
              ? "translateX(-120%) rotate(-14deg)"
              : anim === "right"
              ? "translateX(120%) rotate(14deg)"
              : `translateX(${dragX}px) rotate(${rotation}deg)`,
        }
      : {};

  // Labelit (TRASH/KEEP) läpinäkyvinä pyyhkäisyn mukaan
  const trashOpacity = Math.min(1, Math.max(0, (-dragX - 30) / threshold));
  const keepOpacity = Math.min(1, Math.max(0, (dragX - 30) / threshold));

  return (
    <main className="stage">
      <header className="topbar">
        <h1 className="title">Tiedostot</h1>
        <div className="counter">
          {idx + 1}/{files.length}
        </div>
      </header>

      {/* Pino (renderöi pari seuraavaa varjoksi) */}
      <div className="stack">
        {/* Seuraavat kortit pienellä siirrolla taakse */}
        {files.slice(idx + 1, idx + 3).map((f, i) => (
          <div key={f.path} className={`card shadow-${i + 1}`} aria-hidden />
        ))}

        {/* Aktiivinen kortti */}
        <Card
          file={current}
          style={cardStyle}
          dragging={dragging}
          trashOpacity={trashOpacity}
          keepOpacity={keepOpacity}
          onPointerDown={onPointerDown}
          onPointerMove={onPointerMove}
          onPointerUp={onPointerUp}
          onPointerCancel={onPointerUp}
          anim={anim}
          previewUrl={getPreviewUrl(current)}
        />
      </div>

      <ActionsBar
        onTrash={async () => {
          setAnim("left");
          await new Promise((r) => setTimeout(r, 100));
          await trashCurrent();
        }}
        onKeep={async () => {
          setAnim("right");
          await new Promise((r) => setTimeout(r, 100));
          keepCurrent();
        }}
        remaining={remaining}
      />
    </main>
  );
}