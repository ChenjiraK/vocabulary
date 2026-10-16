import type { ReactNode } from "react";

type MainButtonProps = {
  children: ReactNode;
  onClick?: () => void;
  type?: "button" | "submit" | "reset";
};

export default function MainButton({ children, onClick }: MainButtonProps) {
  return (
    <button
      type="button"
      className="flex h-10 w-36.25 items-center justify-center gap-2 rounded-md border border-[#b9c0c9] bg-primary-blue text-white text-sm font-bold shadow-sm"
      onClick={onClick}
    >
      {children}
    </button>
  );
}
