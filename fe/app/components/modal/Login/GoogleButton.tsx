"use client";

import Image from "next/image";
import { Button } from "@mui/material";
import { loginWithGoogle } from "@/app/services/auth.service";

export default function PasswordButton() {
  async function LoginWithPassword() {
    const data = await loginWithGoogle();
    console.log(data);
  }

  return (
    <div>
      <Button
        variant="outlined"
        onClick={LoginWithPassword}
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
          src="/images/google-icon.png"
          alt=""
          width={24}
          height={24}
          className="mr-3"
        />
        Login with Google
      </Button>
    </div>
  );
}
