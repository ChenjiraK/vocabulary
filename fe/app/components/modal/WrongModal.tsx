"use client";

import { useId } from "react";
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Typography,
} from "@mui/material";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCircleXmark } from "@fortawesome/free-solid-svg-icons";

type WrongModalProps = {
  isOpen: boolean;
  onClose: () => void;
  title?: string;
  message?: string;
  buttonText?: string;
};

export default function WrongModal({
  isOpen,
  onClose,
  title = "Not quite!",
  message = "Sorry, that's not the correct word. Please try again.",
  buttonText = "Close",
}: WrongModalProps) {
  const id = useId();
  const titleId = `${id}-title`;
  const descriptionId = `${id}-description`;

  return (
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
          color: "error.main",
          textAlign: "center",
          fontWeight: 700,
          fontSize: 32,
        }}
      >
        {title}
      </DialogTitle>
      <DialogContent>
        <div className="mb-5 flex justify-center text-[#f02b5d]">
          <FontAwesomeIcon
            icon={faCircleXmark}
            style={{ width: "6rem", height: "6rem" }}
            aria-hidden="true"
          />
        </div>
        <Typography
          id={descriptionId}
          color="text.secondary"
          sx={{ textAlign: "center" }}
        >
          {message}
        </Typography>
      </DialogContent>
      <DialogActions sx={{ justifyContent: "center", m: 2.5 }}>
        <Button
          variant="outlined"
          color="error"
          onClick={onClose}
          autoFocus
          fullWidth
          sx={{ borderRadius: "40px" }}
        >
          {buttonText}
        </Button>
      </DialogActions>
    </Dialog>
  );
}
