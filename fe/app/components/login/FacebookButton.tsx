"use client";

import Image from "next/image";
import { Button } from "@mui/material";
import { loginWithFacebook } from "@/app/services/auth.service";

export default function FacebookButton() {
  async function LoginWithFacebook() {
    const data = await loginWithFacebook();
    console.log(data);
  }

  return (
    <div>
      <Button
        variant="outlined"
        onClick={LoginWithFacebook}
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
