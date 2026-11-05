import Loading from "@/app/components/ui/Loading";

export default function LoadingPage() {
  return (
    <main
      className="grid min-h-screen place-items-center bg-page px-4 py-10"
      aria-busy="true"
    >
      <Loading />
    </main>
  );
}
