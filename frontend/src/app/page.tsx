"use client";
// HomePage: kokoaa sovelluksen tilan (lista + pyyhkäisyt) ja UI:n (Card + ActionsBar).
// Miksi: sivutason tila ja ohjaus yhdessä paikassa.

import { useCallback, useEffect, useRef, useState } from "react";
import "./main.css";
import { fetchFiles, fetchRecents, FileInfo, sendToTrash } from "@/lib/api";
import { Card } from "@/components/Card";
import { ActionsBar } from "@/components/ActionsBar";
import { ChevronDown } from "lucide-react";

type SwipeDir = "left" | "right" | null;

export default function HomePage() {
  // Kansio-valinta. "Äskeiset" = viimeisimmät lataukset (~/Downloads), lajittelu uusimmat ensin.
  const [folder, setFolder] = useState<"Äskeiset" | "Lataukset" | "Työpöytä">("Äskeiset");
  const [menuOpen, setMenuOpen] = useState(false);
  const [files, setFiles] = useState<FileInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Swipen tilat
  const [idx, setIdx] = useState(0);          
  const [dragX, setDragX] = useState(0);     
  const [dragging, setDragging] = useState(false);
  const [anim, setAnim] = useState<SwipeDir>(null);
  const startXRef = useRef<number | null>(null);
  const threshold = 120;              

  // Hae tiedostot valinnan mukaan:
  //  - Äskeiset: useista juurista rekursiivisesti (backend /api/recents)
  //  - Lataukset/Työpöytä: suora listaus yhdestä juuresta
  useEffect(() => {
    setLoading(true);
    setError(null);
    setIdx(0);
    if (folder === "Äskeiset") {
      fetchRecents(300, 2, ["~/Downloads", "~/Desktop", "~/Documents", "~/Pictures"]) 
        .then((list) => setFiles(list))
        .catch((e) => setError(String(e)))
        .finally(() => setLoading(false));
    } else {
      const p = folder === "Lataukset" ? "~/Downloads" : "~/Desktop";
      fetchFiles(p, 500)
        .then((list) => setFiles(list))
        .catch((e) => setError(String(e)))
        .finally(() => setLoading(false));
    }
  }, [folder]);

  const current = files[idx] || null;
  const remaining = Math.max(0, files.length - idx - 1);

  const isImage = useCallback((f: FileInfo | null) => {
    if (!f) return false;
    return /\.(png|jpe?g|gif|webp|bmp|svg)$/i.test(f.ext || f.name);
  }, []);

  // Luo esikatselu-URL kuville backendin /api/open -endpointiin
  const getPreviewUrl = useCallback((f: FileInfo) => {
    const base = (process.env.NEXT_PUBLIC_BACKEND_URL || "").replace(/\/+$/g, "");
    if (!base) return null;
    const target = (f.ext || f.name).toLowerCase();
    const isImg = /\.(png|jpe?g|gif|webp|bmp|svg)$/.test(target);
    const isPdf = /\.(pdf)$/.test(target);
    const isVid = /\.(mp4|webm|ogv|ogg|mov|m4v)$/.test(target);
    const isAud = /\.(mp3|wav|ogg|oga|m4a|aac)$/.test(target);
    const isTxt = /\.(txt|md|markdown|json|log|csv|tsv|js|ts|jsx|tsx|py|go|rs|java|c|cpp|cs|sh|yml|yaml|ini|cfg|toml)$/.test(target);
    const isOffice = /\.(docx|doc|dotx|xlsx|xls|pptx|ppt|odt|odp|ods)$/.test(target);
    if (!(isImg || isPdf || isVid || isAud || isTxt || isOffice)) return null;
    if (isOffice) {
      // Convert office docs to PDF on-the-fly for inline viewing
      return `${base}/api/convert?to=pdf&path=${encodeURIComponent(f.path)}`;
    }
    return `${base}/api/open?path=${encodeURIComponent(f.path)}`;
  }, []);

  // Roskikseen (vasen)
  const trashCurrent = useCallback(async () => {
    const f = files[idx];
    if (!f) return;
    try {
      await sendToTrash(f.path);
      // Poistetaan listasta ja pidetään sama idx (seuraava siirtyy tilalle)
      setFiles((prev) => prev.filter((x, i) => i !== idx));
      setAnim(null);
      setDragX(0);
    } catch (e) {
      alert("Roskikseen siirto epäonnistui: " + e);
    }
  }, [files, idx]);

  // Säilytä (oikea)
  const keepCurrent = useCallback(() => {
    if (idx < files.length) {
      setIdx((i) => Math.min(i + 1, files.length));
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
        <h1 className="title">Ei enempää kortteja. Tiedostot nuoltu puhtaaksi. ✅
        </h1>
        <p className="center">
        Sulje selain tai lataa sivu uudelleen aloittaaksesi alusta.
        </p>
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
        <div className="title folder-select">
          <button
            type="button"
            className="folder-btn"
            onClick={() => setMenuOpen((v) => !v)}
            aria-haspopup="listbox"
            aria-expanded={menuOpen}
            aria-label="Vaihda kansiota"
          >
            {folder}
            <ChevronDown size={18} className="chev" />
          </button>
          {menuOpen && (
            <ul className="folder-menu" role="listbox" aria-label="Kansiot">
              {(["Äskeiset", "Lataukset", "Työpöytä"] as const).map((opt) => (
                <li key={opt}>
                  <button
                    className={`folder-opt ${folder === opt ? "active" : ""}`}
                    role="option"
                    aria-selected={folder === opt}
                    onClick={() => {
                      setFolder(opt);
                      setMenuOpen(false);
                    }}
                  >
                    {opt}
                  </button>
                </li>
              ))}
            </ul>
          )}
        </div>
        <div className="counter">
          {idx + 1}/{files.length}
        </div>
      </header>

      <div className="stack">
        {files.slice(idx + 1, idx + 3).map((f, i) => (
          <div key={f.path} className={`card shadow-${i + 1}`} aria-hidden />
        ))}

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
