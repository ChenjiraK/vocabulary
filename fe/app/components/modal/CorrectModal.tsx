import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
} from "@mui/material";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCircleCheck } from "@fortawesome/free-solid-svg-icons";

type ModalProps = {
  isOpen: boolean;
  onClose?: () => void;
  onSubmit?: () => void;
  title: string;
};
export default function SuccessModal({
  isOpen,
  onClose,
  onSubmit,
  title,
}: ModalProps) {
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
          <div className="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-full bg-green-100 text-green-600">
            <FontAwesomeIcon
              icon={faCircleCheck}
              className="h-24 w-24"
              aria-hidden="true"
            />
          </div>
        </DialogContent>
        <DialogActions sx={{ justifyContent: "right", gap: 1, margin: "20px" }}>
          <Button
            variant="outlined"
            onClick={onClose}
            autoFocus
            sx={{ width: "100%", borderRadius: "40px" }}
            color="success"
          >
            Close
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
}
