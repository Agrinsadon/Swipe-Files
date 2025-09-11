"use client";

import { useEffect, useState } from "react";
import { fetchFiles, FileInfo, sendToTrash } from "@/lib/api";

export default function HomePage() {
  const [files, setFiles] = useState<FileInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchFiles("~/Downloads") // voit vaihtaa polun tähän
        .then(setFiles)
        .catch((e) => setError(String(e)))
        .finally(() => setLoading(false));
  }, []);

  async function handleTrash(path: string) {
    try {
      await sendToTrash(path);
      setFiles((prev) => prev.filter((f) => f.path !== path));
    } catch (e) {
      alert("Roskikseen siirto epäonnistui: " + e);
    }
  }

  if (loading) return <p>Ladataan...</p>;
  if (error) return <p style={{ color: "red" }}>Virhe: {error}</p>;

  return (
      <main className="p-6">
        <h1 className="text-xl font-bold mb-4">Tiedostot</h1>
        <ul className="space-y-2">
          {files.map((f) => (
              <li
                  key={f.path}
                  className="p-2 border rounded flex justify-between items-center"
              >
                <div>
                  <span className="font-mono">{f.name}</span> — {f.size} bytes
                </div>
                <button
                    className="px-3 py-1 bg-red-600 text-white rounded"
                    onClick={() => handleTrash(f.path)}
                >
                  🗑️ Poista
                </button>
              </li>
          ))}
        </ul>
      </main>
  );
}
