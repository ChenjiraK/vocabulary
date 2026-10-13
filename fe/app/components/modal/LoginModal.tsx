"use client";
import { useState } from "react";
import { Dialog, DialogContent } from "@mui/material";
import GoogleButton from "./Login/GoogleButton";
import FacebookButton from "./Login/FacebookButton";
import PasswordButton from "./Login/PasswordButton";
import FormEmailPassword from "./Login/FormEmailPassword";

type ModalProps = {
  isOpen: boolean;
  onClose?: () => void;
};
export default function LoginModal({ isOpen, onClose }: ModalProps) {
  const [isPasswordFormOpen, setIsPasswordFormOpen] = useState(false);
  return (
    <div>
      <Dialog
        open={isOpen}
        onClose={onClose}
        aria-labelledby="submission-result-title"
        aria-describedby="submission-result-description"
        fullWidth
        maxWidth="xs"
      >
        <DialogContent>
          {isPasswordFormOpen ? (
            <div>
              <FormEmailPassword onBack={() => setIsPasswordFormOpen(false)} />
            </div>
          ) : (
            <div className="flex flex-col justify-center py-4">
              <GoogleButton />
              <FacebookButton />
              <PasswordButton onClick={() => setIsPasswordFormOpen(true)} />
            </div>
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
}
