import MainPage from "../components/MainPage";
import Image from "next/image";
import bookIcon from "@/public/images/book-icon.png";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faClock as clockIcon } from "@fortawesome/free-regular-svg-icons/faClock";
import { faStar as starIcon } from "@fortawesome/free-solid-svg-icons/faStar";

const ruleItems = [
  {
    icon: clockIcon,
    iconClassName: "text-[#102b61]",
    text: "Beat the clock! you have 1 minute for each question",
  },
  {
    icon: starIcon,
    iconClassName: "text-[#f5b301]",
    text: "Guess the word correctly to get points!",
  },
];

export default function GamePage() {
  return (
    <MainPage>
      <div className="bg-page flex min-h-screen items-center justify-center px-8 py-10">
        <div className="flex w-full max-w-[320px] flex-col items-center text-center">
          <h1 className="text-[30px] leading-[1.1] font-black tracking-wide text-[#061f59]">
            VOCAB
            <br />
            CHALLENGE
          </h1>

          <div className="mt-8 w-full">
            <Image
              src={bookIcon}
              alt="Book icon"
              width={180}
              height={180}
              className="mx-auto mt-6"
              priority
            />
          </div>

          <div className="mt-9 flex w-full max-w-[270px] flex-col gap-4 text-left text-[15px] text-[#172449]">
            {ruleItems.map((item) => (
              <div key={item.text} className="flex items-center gap-3">
                <span className="flex shrink-0 items-center justify-center">
                  <FontAwesomeIcon
                    icon={item.icon}
                    className={item.iconClassName}
                  />
                </span>
                <span>{item.text}</span>
              </div>
            ))}
          </div>

          <button className="mt-9 h-14 w-full max-w-[300px] rounded-2xl bg-[#102b61] text-xl font-bold text-white shadow-[inset_0_1px_0_rgba(255,255,255,0.2),0_4px_10px_rgba(16,43,97,0.2)] transition hover:bg-[#0b2453]">
            Start
          </button>

          <button className="mt-8 text-base font-bold text-[#102b61]">
            How to play
          </button>
        </div>
      </div>
    </MainPage>
  );
}
