import type { ReactNode } from "react";

type MainPageProps = {
  children: ReactNode;
};

export default function MainPage({ children }: MainPageProps) {
  return (
    <main className="text-text-main min-h-screen bg-page">
      <section className="mx-auto w-full overflow-hidden">
        {children}
      </section>
    </main>
  );
}
