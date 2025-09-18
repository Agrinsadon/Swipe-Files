// FileInfo: backendin palauttama tiedostometa listauksessa.
export type FileInfo = {
    name: string;
    path: string;
    ext: string;
    size: number;
    modTime: string;
};

// fetchFiles: hae hakemiston tiedostot backendistä (uusin ensin).
export async function fetchFiles(dir: string, limit?: number): Promise<FileInfo[]> {
    const base = (process.env.NEXT_PUBLIC_BACKEND_URL || "").replace(/\/+$/g, "");
    if (!base) throw new Error("NEXT_PUBLIC_BACKEND_URL puuttuu");
    const url = `${base}/api/files?dir=${encodeURIComponent(dir)}${
      typeof limit === "number" && limit > 0 ? `&limit=${limit}` : ""
    }`;
    const res = await fetch(url, { cache: "no-store" });
    if (!res.ok) {
        const txt = await res.text().catch(() => "");
        throw new Error(`Virhe ${res.status}: ${txt}`);
    }
    const data = await res.json().catch(() => null);
    return Array.isArray(data) ? (data as FileInfo[]) : [];
}

export async function fetchRecents(limit?: number, maxDepth?: number, roots?: string[]): Promise<FileInfo[]> {
    const base = (process.env.NEXT_PUBLIC_BACKEND_URL || "").replace(/\/+$/g, "");
    if (!base) throw new Error("NEXT_PUBLIC_BACKEND_URL puuttuu");
    const params = new URLSearchParams();
    if (typeof limit === "number" && limit > 0) params.set("limit", String(limit));
    if (typeof maxDepth === "number" && maxDepth >= 0) params.set("maxDepth", String(maxDepth));
    if (roots && roots.length) params.set("dirs", roots.join(","));
    const url = `${base}/api/recents?${params.toString()}`;
    const res = await fetch(url, { cache: "no-store" });
    if (!res.ok) {
        const txt = await res.text().catch(() => "");
        throw new Error(`Virhe ${res.status}: ${txt}`);
    }
    const data = await res.json().catch(() => null);
    return Array.isArray(data) ? (data as FileInfo[]) : [];
}

// sendToTrash: siirrä tiedosto roskakoriin palvelimen kautta.
export async function sendToTrash(path: string): Promise<void> {
    const base = (process.env.NEXT_PUBLIC_BACKEND_URL || "").replace(/\/+$/g, "");
    if (!base) throw new Error("NEXT_PUBLIC_BACKEND_URL puuttuu");
    const res = await fetch(`${base}/api/trash`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ path }),
    });
    if (!res.ok) {
        const txt = await res.text().catch(() => "");
        throw new Error(`Roskikseen siirto epäonnistui ${res.status}: ${txt}`);
    }
}
