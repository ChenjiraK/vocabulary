"use client";

import Image from "next/image";
import { Button } from "@mui/material";

interface PasswordButtonProps {
  onClick?: () => void;
}
export default function PasswordButton({ onClick }: PasswordButtonProps) {
  return (
    <div>
      <Button
        variant="outlined"
        onClick={onClick}
        sx={{
          width: "100%",
          borderRadius: "40px",
          color: "#1f2937",
          borderColor: "#d1d5db",
          textTransform: "none",
          fontWeight: 600,
          py: 1.25,
          my: 1.5,
        }}
      >
        <Image
          src="/images/guest-icon.png"
          alt=""
          width={24}
          height={24}
          className="mr-3"
        />
        Login with Password
      </Button>
    </div>
  );
}
