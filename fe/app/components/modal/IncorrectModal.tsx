import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Typography,
} from "@mui/material";

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
        <DialogTitle id="submission-result-title" sx={{ fontSize: "60px" }}>
          Try again
        </DialogTitle>
        <DialogContent>
          <Typography id="submission-result-description" color="text.secondary">
            Sorry, that's not the correct word. Please try again.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button variant="contained" onClick={onClose} autoFocus>
            Close
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
}
