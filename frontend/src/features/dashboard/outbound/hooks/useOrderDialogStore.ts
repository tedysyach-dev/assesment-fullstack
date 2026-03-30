import { create } from "zustand";

interface OrderDialogState {
  open: boolean;
  orderSn: string | null;
  openDialog: (orderSn: string) => void;
  closeDialog: () => void;
}

export const useOrderDialogStore = create<OrderDialogState>((set) => ({
  open: false,
  orderSn: null,
  openDialog: (orderSn) => {
    document.body.style.overflow = "hidden";
    set({ open: true, orderSn });
  },
  closeDialog: () => {
    document.body.style.overflow = "";
    set({ open: false, orderSn: null });
  },
}));
