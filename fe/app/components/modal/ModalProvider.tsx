"use client";

import {
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
  type ReactNode,
} from "react";
import { usePathname } from "next/navigation";
import { ModalContext, type ModalOptions } from "@/app/hooks/useModal";
import SuccessModal from "./SuccessModal";
import WrongModal from "./WrongModal";

type ModalState = {
  type: "success" | "wrong";
  options: ModalOptions;
  pathname: string;
  isOpen: boolean;
};

export default function ModalProvider({ children }: { children: ReactNode }) {
  const pathname = usePathname();
  const [modal, setModal] = useState<ModalState | null>(null);
  const activeModal = useRef<ModalState | null>(null);

  const dismissModal = useCallback(() => {
    activeModal.current = null;
    setModal((current) => (current ? { ...current, isOpen: false } : null));
  }, []);

  const closeModal = useCallback(() => {
    const current = activeModal.current;
    dismissModal();
    current?.options.onClose?.();
  }, [dismissModal]);

  const showModal = useCallback(
    (type: ModalState["type"], options: ModalOptions = {}) => {
      const nextModal = { type, options, pathname, isOpen: true };
      activeModal.current = nextModal;
      setModal(nextModal);
    },
    [pathname],
  );

  // Dismiss page-specific feedback when navigating, without running its callback.
  useEffect(() => () => dismissModal(), [pathname, dismissModal]);

  const value = useMemo(
    () => ({
      showSuccessModal: (options?: ModalOptions) => showModal("success", options),
      showErrorModal: (options?: ModalOptions) => showModal("wrong", options),
      closeModal,
    }),
    [showModal, closeModal],
  );

  const isOpen = modal?.isOpen === true && modal.pathname === pathname;

  return (
    <ModalContext.Provider value={value}>
      {children}
      {modal?.type === "success" && (
        <SuccessModal {...modal.options} isOpen={isOpen} onClose={closeModal} />
      )}
      {modal?.type === "wrong" && (
        <WrongModal {...modal.options} isOpen={isOpen} onClose={closeModal} />
      )}
    </ModalContext.Provider>
  );
}
