"use client";

interface Props {
  onConfirm: () => void;
  onCancel: () => void;
}

export default function SubmitConfirmModal({ onConfirm, onCancel }: Props) {
  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center"
      style={{ backgroundColor: "rgba(40,40,37,0.5)" }}
      onClick={onCancel}
    >
      <div
        className="bg-white p-8 max-w-sm w-full mx-4"
        style={{
          border: "2px solid rgb(40,40,37)",
          borderRadius: "4px",
          boxShadow: "8px 8px 0 0 rgb(40,40,37)",
        }}
        onClick={(e) => e.stopPropagation()}
      >
        <h2 className="text-xl font-black text-ink mb-2">Поднеси запис?</h2>
        <p className="text-sm opacity-60 leading-relaxed mb-4">
          Овој проект работи на база на доверба. Само преку внесување на точни податоци ќе може да сe добие точна статистика. 🤝
        </p>
        <div className="flex gap-3">
          <button onClick={onConfirm} className="btn-primary flex-1">
            Да, поднеси
          </button>
          <button onClick={onCancel} className="btn-white flex-1">
            Назад
          </button>
        </div>
      </div>
    </div>
  );
}