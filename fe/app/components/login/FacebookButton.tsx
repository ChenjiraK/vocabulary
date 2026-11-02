"use client";

import Image from "next/image";
import { Button } from "@mui/material";
import { loginWithFacebook } from "@/app/services/auth.service";
type FacebookButtonProps = {
  disabled?: boolean;
};

export default function FacebookButton({ disabled }: FacebookButtonProps) {
  async function LoginWithFacebook() {
    const data = await loginWithFacebook();
    console.log(data);
  }

  return (
    <div className={disabled ? "cursor-not-allowed" : undefined}>
      <Button
        variant="outlined"
        onClick={LoginWithFacebook}
        disabled={disabled}
        className={disabled ? "opacity-60" : undefined}
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
          src="/images/facebook-icon.png"
          alt=""
          width={24}
          height={24}
          className="mr-3"
        />
        Login with Facebook
      </Button>
    </div>
  );
}
