"use client";

import { useId } from "react";
import {
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Typography,
} from "@mui/material";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faEnvelopeCircleCheck } from "@fortawesome/free-solid-svg-icons";
import MainButton from "@/app/components/ui/MainButton";

type ModalProps = {
  isOpen: boolean;
  onClose: () => void;
  title?: string;
  message?: string;
  buttonText?: string;
};
export default function SuccessModal({
  isOpen,
  onClose,
  title = "Success!",
  message = "Your action was completed successfully.",
  buttonText = "Close",
}: ModalProps) {
  const id = useId();
  const titleId = `${id}-title`;
  const descriptionId = `${id}-description`;

  return (
    <div>
      <Dialog
        open={isOpen}
        onClose={onClose}
        aria-labelledby={titleId}
        aria-describedby={descriptionId}
        fullWidth
        maxWidth="xs"
      >
        <DialogTitle
          id={titleId}
          sx={{
            color: "success.main",
            textAlign: "center",
            fontWeight: 700,
            fontSize: 32,
          }}
        >
          {title}
        </DialogTitle>
        <DialogContent>
          <div className="flex flex-col items-center justify-center">
            <span className="text-success mt-3">
              <FontAwesomeIcon
                icon={faEnvelopeCircleCheck}
                className="text-9xl"
                aria-hidden="true"
              />
            </span>
          </div>
          <Typography
            id={descriptionId}
            color="text.secondary"
            className="pt-7 text-center text-lg font-semibold"
          >
            {message}
          </Typography>
        </DialogContent>
        <DialogActions sx={{ justifyContent: "center", p: 2.5 }}>
          <MainButton onClick={onClose}>{buttonText}</MainButton>
        </DialogActions>
      </Dialog>
    </div>
  );
}
