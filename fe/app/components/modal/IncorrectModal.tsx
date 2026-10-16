import {
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Typography,
} from "@mui/material";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCircleXmark } from "@fortawesome/free-solid-svg-icons";
import MainButton from "@/app/components/ui/MainButton";

type ModalProps = {
  isOpen: boolean;
  onClose: () => void;
};
export default function IncorrectModal({ isOpen, onClose }: ModalProps) {
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
          <div className="flex flex-col items-center justify-center">
            <span className="text-[#f02b5d] mt-3">
              <FontAwesomeIcon icon={faCircleXmark} className="text-9xl" />
            </span>
          </div>
          <Typography
            id="submission-result-description"
            color="text.secondary"
            className="pt-7 text-center text-lg font-semibold"
          >
            Sorry, that's not the correct word. Please try again.
          </Typography>
        </DialogContent>
        <DialogActions>
          <MainButton onClick={onClose}>Try Again</MainButton>
        </DialogActions>
      </Dialog>
    </div>
  );
}
