// @/lib/date.ts
import { format, parseISO } from "date-fns";
import { id } from "date-fns/locale";

export const formatDate = (isoString: string): string => {
  return format(parseISO(isoString), "dd MMM yyyy", { locale: id });
};

export const formatDateTime = (isoString: string): string => {
  return format(parseISO(isoString), "dd MMM yyyy HH:mm", { locale: id });
};

export const formatDateTimeSecond = (isoString: string): string => {
  return format(parseISO(isoString), "dd MMM yyyy HH:mm:ss", { locale: id });
};
