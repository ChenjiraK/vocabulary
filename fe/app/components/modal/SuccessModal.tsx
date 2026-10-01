import Image from "next/image";
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
} from "@mui/material";

type ModalProps = {
  isOpen: boolean;
  onClose?: () => void;
  onSubmit?: () => void;
  title: string;
};
export default function SuccessModal({ isOpen, onClose,onSubmit, title }: ModalProps) {
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
        <DialogTitle
          id="submission-result-title"
          className="text-success text-center font-bold"
          sx={{ fontSize: "32px" }}
        >
          {title || "Congratulations!"}
        </DialogTitle>
        <DialogContent>
          <div className="flex justify-center">
            <Image
              src="/images/star-icon.gif"
              alt="Correct answer"
              width={200}
              height={200}
              unoptimized
            />
          </div>
        </DialogContent>
        <DialogActions sx={{ justifyContent: "center", gap: 1, margin: "20px" }}>
          <Button variant="contained" onClick={onSubmit} autoFocus sx={{ width: "100%", borderRadius:"40px" }} color="success">
            Continue
          </Button>
          <Button variant="outlined" onClick={onClose} autoFocus sx={{ width: "100%", borderRadius:"40px" }} color="success">
            Close
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
}
