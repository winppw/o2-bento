export interface WindowStatus {
  isOpen: boolean;
  openAt: Date;
  closeAt: Date;
  isFridayExtra: boolean;
  message: string;
}

export function getWindowStatus(now: Date = new Date()): WindowStatus {
  // Convert to Asia/Bangkok (UTC+7)
  const bkk = new Date(now.toLocaleString("en-US", { timeZone: "Asia/Bangkok" }));

  const open = new Date(bkk);
  open.setHours(11, 0, 0, 0);

  const close = new Date(bkk);
  close.setHours(16, 0, 0, 0);

  const isWeekend = bkk.getDay() === 0 || bkk.getDay() === 6;
  const isOpen = !isWeekend && bkk >= open && bkk < close;
  const isFriday = bkk.getDay() === 5;

  let message = "";
  if (bkk < open) {
    message = "Order window opens at 11:00 AM.";
  } else if (isOpen && isFriday) {
    message = "Order window is open. You can also submit Monday's order today.";
  } else if (isOpen) {
    message = "Order window is open. Deadline is 4:00 PM.";
  } else {
    message = "Order window is closed for today.";
  }

  return { isOpen, openAt: open, closeAt: close, isFridayExtra: isFriday, message };
}

export function formatTimeRemaining(target: Date, now: Date = new Date()): string {
  const diff = target.getTime() - now.getTime();
  if (diff <= 0) return "0:00";
  const h = Math.floor(diff / 3_600_000);
  const m = Math.floor((diff % 3_600_000) / 60_000);
  const s = Math.floor((diff % 60_000) / 1_000);
  if (h > 0) return `${h}h ${m}m`;
  return `${m}:${s.toString().padStart(2, "0")}`;
}
