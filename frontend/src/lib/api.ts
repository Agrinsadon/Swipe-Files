export type FileInfo = {
    name: string;
    path: string;
    ext: string;
    size: number;
    modTime: string;
};

export async function fetchFiles(dir: string): Promise<FileInfo[]> {
    const base = (process.env.NEXT_PUBLIC_BACKEND_URL || "").replace(/\/+$/g, "");
    if (!base) throw new Error("NEXT_PUBLIC_BACKEND_URL puuttuu");
    const url = `${base}/api/files?dir=${encodeURIComponent(dir)}`;
    const res = await fetch(url, { cache: "no-store" });
    if (!res.ok) {
        const txt = await res.text().catch(() => "");
        throw new Error(`Error ${res.status}: ${txt}`);
    }
    return res.json();
}

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
        throw new Error(`Trash failed ${res.status}: ${txt}`);
    }
}
