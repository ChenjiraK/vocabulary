"use client";

import { Dialog } from "@mui/material";
import Loading, { type LoadingProps } from "@/app/components/ui/Loading";

type LoadingModalProps = LoadingProps & {
  isOpen: boolean;
};

export default function LoadingModal({
  isOpen,
  title = "Just a moment…",
  message,
}: LoadingModalProps) {
  return (
    <Dialog
      open={isOpen}
      aria-label={title}
      fullWidth
      maxWidth="xs"
      slotProps={{
        backdrop: {
          sx: {
            backgroundColor: "rgba(9, 26, 53, 0.4)",
            backdropFilter: "blur(4px)",
          },
        },
        paper: {
          sx: {
            m: 2,
            width: "calc(100% - 32px)",
            maxWidth: 384,
            borderRadius: "28px",
            backgroundColor: "#fffdf8",
            boxShadow: "0 24px 80px rgba(9, 26, 53, 0.2)",
          },
        },
      }}
    >
      <Loading title={title} message={message} />
    </Dialog>
  );
}
