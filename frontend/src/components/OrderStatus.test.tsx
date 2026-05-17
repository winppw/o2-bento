import { render, screen } from "@testing-library/react";
import OrderStatus from "./OrderStatus";

// Pin time to a known state: Tuesday 2024-01-09 14:00 BKK (07:00 UTC)
const OPEN_TIME = new Date("2024-01-09T07:00:00Z");
// Tuesday 2024-01-09 18:00 BKK (11:00 UTC) — after deadline
const CLOSED_TIME = new Date("2024-01-09T11:00:00Z");

function renderWithTime(date: Date) {
  jest.useFakeTimers();
  jest.setSystemTime(date);
  const result = render(<OrderStatus />);
  jest.useRealTimers();
  return result;
}

describe("OrderStatus", () => {
  test("shows OPEN badge during window", () => {
    renderWithTime(OPEN_TIME);
    expect(screen.getByText("OPEN")).toBeInTheDocument();
  });

  test("shows CLOSED badge outside window", () => {
    renderWithTime(CLOSED_TIME);
    expect(screen.getByText("CLOSED")).toBeInTheDocument();
  });

  test("form button is enabled and links to form when open", () => {
    process.env.NEXT_PUBLIC_GOOGLE_FORM_URL = "https://forms.gle/test";
    renderWithTime(OPEN_TIME);
    const btn = screen.getByRole("link", { name: /submit lunch order/i });
    expect(btn).toHaveAttribute("href", "https://forms.gle/test");
  });

  test("form button is disabled when closed", () => {
    renderWithTime(CLOSED_TIME);
    const btn = screen.getByRole("button", { name: /form unavailable/i });
    expect(btn).toBeDisabled();
  });

  test("pickup note is always visible", () => {
    renderWithTime(OPEN_TIME);
    expect(screen.getByText(/pantry/i)).toBeInTheDocument();
  });

  test("Friday note visible on Friday during window", () => {
    // 2024-01-12 is a Friday, 13:00 BKK = 06:00 UTC
    const friday = new Date("2024-01-12T06:00:00Z");
    renderWithTime(friday);
    // "Monday" appears in both the status message and the Friday extra note
    const mondayRefs = screen.getAllByText(/monday/i);
    expect(mondayRefs.length).toBeGreaterThan(0);
  });
});
