import { useEffect, useState } from "react";

export default function CircleCountdown() {
  const timerRadius = 22;
  const countdownSeconds = 60;
  const timerCircumference = 2 * Math.PI * timerRadius;
  const [timeLeft, setTimeLeft] = useState(countdownSeconds);
  const displayTime = formatTime(timeLeft);
  const timerProgress = timeLeft / countdownSeconds;
  const timerDashOffset = timerCircumference * (1 - timerProgress);

  function formatTime(totalSeconds: number) {
    const minutes = Math.floor(totalSeconds / 60);
    const seconds = totalSeconds % 60;

    return `${String(minutes).padStart(2, "0")}:${String(seconds).padStart(2, "0")}`;
  }

  function countdownTimer() {
    //countdown timer
    const timerId = window.setInterval(() => {
      setTimeLeft((currentTimeLeft) => {
        if (currentTimeLeft <= 1) {
          window.clearInterval(timerId);

          return 0;
        }

        return currentTimeLeft - 1;
      });
    }, 1000);

    return () => window.clearInterval(timerId);
  }

  useEffect(() => {
    setTimeLeft(countdownSeconds);
    return countdownTimer();
  }, []);

  return (
    <div
      id="timer"
      className="relative flex h-12 w-12 items-center justify-center rounded-full bg-white text-[11px] font-semibold text-slate-950"
    >
      <svg
        viewBox="0 0 48 48"
        className="absolute inset-0 h-full w-full -rotate-90"
        aria-hidden="true"
      >
        <circle
          cx="24"
          cy="24"
          r={timerRadius}
          fill="none"
          stroke="#dbe3f0"
          strokeWidth="3"
        />
        <circle
          cx="24"
          cy="24"
          r={timerRadius}
          fill="none"
          stroke="#003dcc"
          strokeWidth="3"
          strokeLinecap="round"
          strokeDasharray={timerCircumference}
          strokeDashoffset={timerDashOffset}
          className="transition-[stroke-dashoffset] duration-1000 ease-linear"
        />
      </svg>
      {displayTime}
    </div>
  );
}
