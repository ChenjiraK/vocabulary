import { redirect } from "next/navigation";

export default function Home() {
  //set welcome to first page
  redirect("/welcome");
}
