import { getWindowStatus, formatTimeRemaining } from "./schedule";

// Helper: build a Date that appears as the given BKK local time.
// We simulate by constructing a UTC time that, when converted to Asia/Bangkok (UTC+7),
// reads as the intended h:m on a known weekday.
function bkk(isoDate: string, hour: number, minute: number): Date {
  // isoDate like "2024-01-09" (a Tuesday)
  const [y, m, d] = isoDate.split("-").map(Number);
  // Asia/Bangkok is UTC+7, so UTC hour = BKK hour - 7
  return new Date(Date.UTC(y, m - 1, d, hour - 7, minute, 0));
}

describe("getWindowStatus", () => {
  test("before open returns isOpen=false", () => {
    // Tuesday 10:59 AM BKK
    const { isOpen, message } = getWindowStatus(bkk("2024-01-09", 10, 59));
    expect(isOpen).toBe(false);
    expect(message).toMatch(/opens at 11/i);
  });

  test("at 11:00 AM returns isOpen=true", () => {
    const { isOpen } = getWindowStatus(bkk("2024-01-09", 11, 0));
    expect(isOpen).toBe(true);
  });

  test("at 3:59 PM returns isOpen=true", () => {
    const { isOpen } = getWindowStatus(bkk("2024-01-09", 15, 59));
    expect(isOpen).toBe(true);
  });

  test("at 4:00 PM returns isOpen=false", () => {
    const { isOpen, message } = getWindowStatus(bkk("2024-01-09", 16, 0));
    expect(isOpen).toBe(false);
    expect(message).toMatch(/closed/i);
  });

  test("after 4:00 PM returns isOpen=false", () => {
    const { isOpen } = getWindowStatus(bkk("2024-01-09", 17, 30));
    expect(isOpen).toBe(false);
  });

  test("Friday open sets isFridayExtra=true and message mentions Monday", () => {
    // 2024-01-12 is a Friday
    const { isOpen, isFridayExtra, message } = getWindowStatus(bkk("2024-01-12", 13, 0));
    expect(isOpen).toBe(true);
    expect(isFridayExtra).toBe(true);
    expect(message).toMatch(/monday/i);
  });

  test("Friday before open: isFridayExtra=true but isOpen=false", () => {
    const { isOpen, isFridayExtra } = getWindowStatus(bkk("2024-01-12", 9, 0));
    expect(isOpen).toBe(false);
    expect(isFridayExtra).toBe(true);
  });

  test("Saturday is closed even during window hours", () => {
    // 2024-01-13 is Saturday
    const { isOpen } = getWindowStatus(bkk("2024-01-13", 13, 0));
    expect(isOpen).toBe(false);
  });

  test("openAt is 11:00 in BKK for every call", () => {
    const { openAt } = getWindowStatus(bkk("2024-01-09", 14, 0));
    const bkkOpen = new Date(openAt.toLocaleString("en-US", { timeZone: "Asia/Bangkok" }));
    expect(bkkOpen.getHours()).toBe(11);
    expect(bkkOpen.getMinutes()).toBe(0);
  });

  test("closeAt is 16:00 in BKK for every call", () => {
    const { closeAt } = getWindowStatus(bkk("2024-01-09", 14, 0));
    const bkkClose = new Date(closeAt.toLocaleString("en-US", { timeZone: "Asia/Bangkok" }));
    expect(bkkClose.getHours()).toBe(16);
    expect(bkkClose.getMinutes()).toBe(0);
  });
});

describe("formatTimeRemaining", () => {
  const base = new Date("2024-01-09T07:00:00Z"); // 2:00 PM BKK

  test("returns h m format for > 60 minutes", () => {
    const target = new Date(base.getTime() + 90 * 60 * 1000);
    expect(formatTimeRemaining(target, base)).toBe("1h 30m");
  });

  test("returns m:ss format for < 60 minutes", () => {
    const target = new Date(base.getTime() + 5 * 60 * 1000 + 3 * 1000);
    expect(formatTimeRemaining(target, base)).toBe("5:03");
  });

  test("returns 0:00 for past target", () => {
    const target = new Date(base.getTime() - 1000);
    expect(formatTimeRemaining(target, base)).toBe("0:00");
  });

  test("pads seconds to two digits", () => {
    const target = new Date(base.getTime() + 65 * 1000);
    expect(formatTimeRemaining(target, base)).toBe("1:05");
  });
});
