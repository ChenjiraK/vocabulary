"use client";

import { createContext, useContext } from "react";

export type ModalOptions = {
  title?: string;
  message?: string;
  buttonText?: string;
  onClose?: () => void;
};

export type LoadingModalOptions = Pick<ModalOptions, "title" | "message">;

type ModalContextValue = {
  showSuccessModal: (options?: ModalOptions) => void;
  showErrorModal: (options?: ModalOptions) => void;
  showLoadingModal: (options?: LoadingModalOptions) => void;
  hideLoadingModal: () => void;
  isLoading: boolean;
  closeModal: () => void;
};

export const ModalContext = createContext<ModalContextValue | null>(null);

export function useModal() {
  const context = useContext(ModalContext);

  if (!context) {
    throw new Error("useModal must be used within ModalProvider.");
  }

  return context;
}
