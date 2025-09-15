"use client";

type Props = {
  onTrash: () => Promise<void> | void;
  onKeep: () => Promise<void> | void;
  remaining: number;
};

export function ActionsBar({ onTrash, onKeep, remaining }: Props) {
  return (
    <>
      <nav className="actions">
        <button
          className="btn danger"
          onClick={async () => {
            await onTrash();
          }}
          aria-label="Siirrä roskakoriin (vasen nuoli)"
        >
          ← Poista
        </button>
        <button
          className="btn"
          onClick={async () => {
            await onKeep();
          }}
          aria-label="Säilytä ja siirry seuraavaan (oikea nuoli)"
        >
          Säilytä →
        </button>
      </nav>
      <p className="hint">Vinkki: käytä nuolia ← / → tai raahaa korttia. {remaining} jäljellä.</p>
    </>
  );
}

