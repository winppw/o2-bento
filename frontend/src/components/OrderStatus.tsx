"use client";

import { useEffect, useState } from "react";
import { getWindowStatus, formatTimeRemaining } from "@/lib/schedule";

export default function OrderStatus() {
  const [now, setNow] = useState(new Date());
  const formURL = process.env.NEXT_PUBLIC_GOOGLE_FORM_URL ?? "#";

  useEffect(() => {
    const id = setInterval(() => setNow(new Date()), 1000);
    return () => clearInterval(id);
  }, []);

  const { isOpen, openAt, closeAt, isFridayExtra, message } = getWindowStatus(now);
  const target = isOpen ? closeAt : openAt < now ? openAt : openAt;
  const countdown = isOpen
    ? formatTimeRemaining(closeAt, now)
    : now < openAt
    ? formatTimeRemaining(openAt, now)
    : null;

  return (
    <div className="flex flex-col items-center gap-6 p-8 max-w-md mx-auto">
      {/* Status badge */}
      <div
        className={`px-4 py-1 rounded-full text-sm font-semibold ${
          isOpen ? "bg-green-100 text-green-700" : "bg-red-100 text-red-600"
        }`}
      >
        {isOpen ? "OPEN" : "CLOSED"}
      </div>

      {/* Message */}
      <p className="text-center text-gray-700">{message}</p>

      {/* Friday note */}
      {isFridayExtra && isOpen && (
        <p className="text-sm text-blue-600 font-medium">
          📅 You can also submit your Monday order today.
        </p>
      )}

      {/* Countdown */}
      {countdown && (
        <div className="text-4xl font-mono font-bold text-gray-800">
          {countdown}
          <span className="text-base font-normal text-gray-500 ml-2">
            {isOpen ? "until close" : "until open"}
          </span>
        </div>
      )}

      {/* Form button */}
      {isOpen ? (
        <a
          href={formURL}
          target="_blank"
          rel="noopener noreferrer"
          className="w-full text-center bg-green-600 hover:bg-green-700 text-white font-semibold py-3 px-6 rounded-xl transition"
        >
          Submit Lunch Order →
        </a>
      ) : (
        <button
          disabled
          className="w-full text-center bg-gray-200 text-gray-400 font-semibold py-3 px-6 rounded-xl cursor-not-allowed"
        >
          Form Unavailable
        </button>
      )}

      {/* Pickup note */}
      <p className="text-xs text-gray-400 text-center">
        Pickup at the pantry: <strong>11:00 AM – 12:00 PM</strong> the following day.
        Check the paper sheet for your order number.
      </p>
    </div>
  );
}
