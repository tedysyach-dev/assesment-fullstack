// pages/OrderPage.tsx
import { useState, useMemo, useEffect } from "react";
import { formatDateTimeSecond } from "@/lib/date";
import { useOrders } from "../hooks/useOrders";
import type { Order } from "../types";
import {
  useReactTable,
  getCoreRowModel,
  getSortedRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  flexRender,
  type ColumnDef,
  type SortingState,
  type ColumnFiltersState,
  type PaginationState,
} from "@tanstack/react-table";
import {
  ChevronsUpDown,
  ChevronDown,
  ChevronUp,
  Search,
  ChevronRight,
  ChevronLeft,
} from "lucide-react";
import ActionOrderButton from "./ActionButton";
import StatsCard from "./StatsCard";
import { Separator } from "@/components/seperator";

const OrderPage = () => {
  const { data: orders, isLoading, isError } = useOrders();

  const [sorting, setSorting] = useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([]);
  const [globalFilter, setGlobalFilter] = useState("");
  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize: 10,
  });

  const columns = useMemo<ColumnDef<Order>[]>(
    () => [
      {
        accessorKey: "order_sn",
        header: "Order SN",
        enableSorting: false,
        cell: (info) => (
          <span className="text-xss text-neutral-100">
            {info.getValue() as string}
          </span>
        ),
      },
      {
        accessorKey: "marketplace_status",
        header: "Marketplace Status",
        meta: {
          className: "text-center",
        },
        cell: (info) => {
          const value = info.getValue() as string;
          const colorMap: Record<string, string> = {
            PAID: "bg-info-surface text-info-main",
            SHIPPING: "bg-warning-surface text-warning-main",
            PROCESSING: "bg-process-surface text-process-main",
            CANCELLED: "bg-danger-surface text-danger-main",
            DELIVERED: "bg-success-surface text-success-main",
          };
          const colorClass =
            colorMap[value?.toUpperCase()] ?? "bg-gray-100 text-gray-700";
          return (
            <span className={`px-2 py-1 rounded-sm text-xss ${colorClass}`}>
              {value
                ?.toLowerCase()
                .replaceAll("_", " ")
                .replace(/\b\w/g, (c) => c.toUpperCase())}
            </span>
          );
        },
      },
      {
        accessorKey: "shipping_status",
        header: "Shipping Status",
        meta: {
          className: "text-center",
        },
        cell: (info) => {
          const value = info.getValue() as string;
          const colorMap: Record<string, string> = {
            SHIPPED: "bg-indigo-100 text-indigo-700",
            LABEL_CREATED: "bg-info-surface text-info-main",
            DELIVERED: "bg-green-100 text-green-700",
            AWAITING_PICKUP: "bg-process-surface text-process-main",
            CANCELLED: "bg-danger-surface text-danger-main",
          };
          const colorClass =
            colorMap[value?.toUpperCase()] ?? "bg-gray-100 text-gray-700";
          return (
            <span className={`px-2 py-1 rounded-sm text-xss ${colorClass}`}>
              {value
                ?.toLowerCase()
                .replaceAll("_", " ")
                .replace(/\b\w/g, (c) => c.toUpperCase())}
            </span>
          );
        },
      },
      {
        accessorKey: "wms_status",
        header: "WMS Status",
        meta: {
          className: "text-center",
        },
        cell: (info) => {
          const value = info.getValue() as string;
          const colorMap: Record<string, string> = {
            READY_TO_PICK: "bg-warning-surface text-warning-main",
            PICKING: "bg-info-surface text-info-main",
            PACKED: "bg-process-surface text-process-main",
            SHIPED: "bg-success-surface text-success-main",
          };
          const colorClass =
            colorMap[value?.toUpperCase()] ?? "bg-gray-100 text-gray-700";
          return (
            <span className={`px-2 py-1 rounded-sm text-xss ${colorClass}`}>
              {value
                ?.toLowerCase()
                .replaceAll("_", " ")
                .replace(/\b\w/g, (c) => c.toUpperCase())}
            </span>
          );
        },
      },
      {
        accessorKey: "tracking_number",
        header: "Tracking Number",
        enableSorting: false,
        cell: (info) => {
          const value = info.getValue() as string;
          return <span className={`px-2 py-1 text-xss`}>{value}</span>;
        },
      },
      {
        accessorKey: "updated_at",
        header: "Updated At",
        cell: (info) => (
          <span className="text-xss text-gray-500">
            {formatDateTimeSecond(info.getValue() as string)}
          </span>
        ),
      },
      {
        id: "actions",
        header: "Actions",
        enableSorting: false,
        cell: ({ row }) => (
          <ActionOrderButton orderSn={row.original.order_sn} />
        ),
      },
    ],
    [],
  );

  const table = useReactTable({
    data: orders ?? [],
    columns,
    state: {
      sorting,
      columnFilters,
      globalFilter,
      pagination,
    },
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    onGlobalFilterChange: setGlobalFilter,
    onPaginationChange: setPagination,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getPaginationRowModel: getPaginationRowModel(),

    autoResetPageIndex: false,
  });

  const currentPage = table.getState().pagination.pageIndex;
  const pageCount = table.getPageCount();

  const stats = useMemo(() => {
    let total = 0;
    let cancelled = 0;

    orders?.forEach((o) => {
      total++;

      if (o.marketplace_status === "cancelled") cancelled++;
    });

    return {
      total,
      cancelled,
    };
  }, [orders]);

  useEffect(() => {
    if (pagination.pageIndex > table.getPageCount() - 1) {
      setPagination((prev) => ({
        ...prev,
        pageIndex: Math.max(0, table.getPageCount() - 1),
      }));
    }
  }, [orders]);

  const getPages = () => {
    const delta = 2;
    const pages: (number | "...")[] = [];

    const rangeStart = Math.max(0, currentPage - delta);
    const rangeEnd = Math.min(pageCount - 1, currentPage + delta);

    if (rangeStart > 0) {
      pages.push(0);
      if (rangeStart > 1) pages.push("...");
    }

    for (let i = rangeStart; i <= rangeEnd; i++) {
      pages.push(i);
    }

    if (rangeEnd < pageCount - 1) {
      if (rangeEnd < pageCount - 2) pages.push("...");
      pages.push(pageCount - 1);
    }

    return pages;
  };

  const pages = getPages();

  if (isLoading)
    return (
      <div className="flex items-center justify-center h-64 text-gray-500">
        Loading...
      </div>
    );

  if (isError)
    return (
      <div className="flex items-center justify-center h-64 text-red-500">
        Gagal memuat data
      </div>
    );

  return (
    <div className="p-6 space-y-4 flex flex-col gap-8">
      <div className="flex flex-col gap-2">
        <h1 className="text-bm font-bold text-neutral-100">Outbound</h1>
        <p className="text-xs font-light text-neutral-80">
          Manage all outbound proccess
        </p>
      </div>
      <div className="flex flex-row gap-5">
        <StatsCard
          title="Total Order"
          value={stats.total}
          status={true}
          statusValue={15}
        />
        <StatsCard
          title="Cancelled"
          value={stats.cancelled}
          status={false}
          statusValue={15}
        />
      </div>
      <div className="w-full p-2 border flex flex-row gap-4 rounded-lg border-neutral-50">
        <div className="relative w-126.25">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-black" />
          <input
            type="text"
            value={globalFilter}
            onChange={(e) => setGlobalFilter(e.target.value)}
            placeholder="Search here..."
            className="w-full border border-neutral-50 rounded-md pl-9 pr-4 py-2 text-sm
               shadow-[inset_1px_2px_2px_0_rgba(0,0,0,0.08)]
               focus:outline-none focus:ring-2 focus:ring-primary-main"
          />
        </div>
        <Separator orientation="vertical" />
      </div>

      {/* Table */}
      <div className="flex flex-col gap-5.75">
        <div className="overflow-x-auto rounded-xl border border-neutral-50 shadow-sm">
          <table className="min-w-full divide-y divide-neutral-50 bg-white">
            <thead className="bg-white">
              {table.getHeaderGroups().map((headerGroup) => (
                <tr key={headerGroup.id}>
                  {headerGroup.headers.map((header) => (
                    <th
                      key={header.id}
                      onClick={header.column.getToggleSortingHandler()}
                      className={`px-4 py-3.5 text-left text-xs font-normal text-gray-500 tracking-wider select-none ${
                        header.column.getCanSort()
                          ? "cursor-pointer hover:text-gray-800"
                          : ""
                      }`}
                    >
                      <div className="flex flex-row  items-center justify-between">
                        {flexRender(
                          header.column.columnDef.header,
                          header.getContext(),
                        )}
                        {header.column.getCanSort() && (
                          <span className="text-black">
                            {header.column.getIsSorted() === "asc" ? (
                              <ChevronUp className="w-3.5 h-3.5 text-black" />
                            ) : header.column.getIsSorted() === "desc" ? (
                              <ChevronDown className="w-3.5 h-3.5 text-black" />
                            ) : (
                              <ChevronsUpDown className="w-3.5 h-3.5 text-black" />
                            )}
                          </span>
                        )}
                      </div>
                    </th>
                  ))}
                </tr>
              ))}
            </thead>
            <tbody className="divide-y divide-gray-100">
              {table.getRowModel().rows.length === 0 ? (
                <tr>
                  <td
                    colSpan={columns.length}
                    className="text-center py-10 text-gray-400 text-sm"
                  >
                    Tidak ada data order
                  </td>
                </tr>
              ) : (
                table.getRowModel().rows.map((row) => (
                  <tr
                    key={row.id}
                    className="odd:bg-white even:bg-neutral-20 hover:bg-blue-50 transition-colors"
                  >
                    {row.getVisibleCells().map((cell) => (
                      <td
                        key={cell.id}
                        className={`px-4 py-3 ${cell.column.columnDef.meta?.className ?? ""}`}
                      >
                        {cell.column.id === "actions" || cell.getValue()
                          ? flexRender(
                              cell.column.columnDef.cell,
                              cell.getContext(),
                            )
                          : "-"}
                      </td>
                    ))}
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        <div className="flex items-center justify-between text-xs text-gray-600">
          <span>
            Show{" "}
            {table.getState().pagination.pageIndex *
              table.getState().pagination.pageSize +
              1}{" "}
            to{" "}
            <select
              value={table.getState().pagination.pageSize}
              onChange={(e) => table.setPageSize(Number(e.target.value))}
              className="border border-gray-300 rounded text-xs"
            >
              {[10, 20, 50].map((size) => (
                <option key={size} value={size}>
                  {size}
                </option>
              ))}
            </select>{" "}
            of {table.getFilteredRowModel().rows.length} entries
          </span>

          <div className="flex items-center gap-1">
            <button
              onClick={() => table.previousPage()}
              disabled={!table.getCanPreviousPage()}
              className="w-7 h-7 border border-neutral-50 rounded disabled:opacity-40"
            >
              <ChevronLeft className="h-4" />
            </button>

            {pages.map((page, i) =>
              page === "..." ? (
                <span key={i} className="px-2">
                  ...
                </span>
              ) : (
                <button
                  key={i}
                  onClick={() => table.setPageIndex(page)}
                  className={`h-7 w-7 border rounded
        ${
          currentPage === page
            ? "bg-primary-surface text-neutral-100 border-primary-main"
            : "hover:bg-gray-100"
        }`}
                >
                  {page + 1}
                </button>
              ),
            )}

            <button
              onClick={() => table.nextPage()}
              disabled={!table.getCanNextPage()}
              className="w-7 h-7 border border-neutral-50 rounded disabled:opacity-40"
            >
              <ChevronRight className="h-4" />
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default OrderPage;
