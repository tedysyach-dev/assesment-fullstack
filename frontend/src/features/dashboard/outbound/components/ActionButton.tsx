// components/ActionOrderButton.tsx

import { Button } from "@/components/button";
import { useOrderDialogStore } from "../hooks/useOrderDialogStore";
import {
  useOrderDetail,
  usePackOrder,
  usePickOrder,
  useShipOrder,
} from "../hooks/useOrders";
import { formatDateTimeSecond } from "@/lib/date";
import {
  useReactTable,
  getCoreRowModel,
  flexRender,
  type ColumnDef,
} from "@tanstack/react-table";
import { useMemo } from "react";
import type { OrderItem } from "../types";

interface ActionOrderButtonProps {
  orderSn: string;
}
const ActionOrderButton = ({ orderSn }: ActionOrderButtonProps) => {
  const {
    open,
    orderSn: activeOrderSn,
    openDialog,
    closeDialog,
  } = useOrderDialogStore();

  const isCurrentOrder = activeOrderSn === orderSn;

  const { data, isLoading } = useOrderDetail(activeOrderSn ?? "");

  const { mutate: pickOrder, isPending: isPicking } = usePickOrder();
  const { mutate: packOrder, isPending: isPacking } = usePackOrder();
  const { mutate: shipOrder, isPending: isShipping } = useShipOrder();

  const actionConfig = {
    READY_TO_PICK: {
      label: "Pickup",
      action: () => pickOrder(orderSn),
      loading: isPicking,
    },

    PICKING: {
      label: "Pack",
      action: () => packOrder(orderSn),
      loading: isPacking,
    },

    PACKED: {
      label: "Ship",
      action: () =>
        shipOrder({
          orderSn: orderSn,
          channelId: "JNE",
        }),
      loading: isShipping,
    },
  } as const;

  const currentAction = data?.wms_status
    ? actionConfig[data.wms_status as keyof typeof actionConfig]
    : null;

  const itemColumns = useMemo<ColumnDef<OrderItem>[]>(
    () => [
      {
        accessorKey: "sku",
        header: "SKU",
        cell: (info) => (
          <span className="text-xss text-neutral-100">
            {info.getValue() as string}
          </span>
        ),
      },
      {
        accessorKey: "quantity",
        header: "QTY",
        cell: (info) => (
          <span className="text-xss  text-neutral-100">
            {info.getValue() as number}
          </span>
        ),
      },
      {
        accessorKey: "price",
        header: "Price",
        cell: (info) => (
          <span className="text-xss text-neutral-100">
            Rp {(info.getValue() as number).toLocaleString("id-ID")}
          </span>
        ),
      },
    ],
    [],
  );

  const table = useReactTable({
    data: data?.items ?? [],
    columns: itemColumns,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <>
      <Button
        variant={"default"}
        size={"default"}
        onClick={() => openDialog(orderSn)}
      >
        Detail
      </Button>

      {open && isCurrentOrder && (
        <>
          {/* backdrop */}
          <div
            className="fixed inset-0 bg-black/30 z-40"
            onClick={closeDialog}
          />

          {/* dialog */}
          <div className="fixed left-1/2 top-1/2 z-50 w-98.25 max-w-lg -translate-x-1/2 -translate-y-1/2 bg-white rounded-2xl shadow-xl p-4 space-y-4">
            {/* header */}
            <div className="flex justify-between items-center">
              <h2 className="text-sm font-medium">Detail</h2>
            </div>

            {/* content */}
            {isLoading ? (
              <p className="text-sm text-gray-400">Loading...</p>
            ) : !data ? (
              <p className="text-sm text-red-500">Data tidak ditemukan</p>
            ) : (
              <>
                {/* info grid */}
                <div className="grid grid-cols-2 gap-x-8 gap-y-4 text-sm">
                  <div>
                    <p className="text-neutral-80 text-xss">Order SN</p>
                    <p className="font-medium text-neutral-100 text-xs">
                      {data.order_sn}
                    </p>
                  </div>

                  <div>
                    <p className="text-neutral-80 text-xss">Order SN</p>
                    <p className="font-medium text-neutral-100 text-xs">
                      {data.order_sn}
                    </p>
                  </div>

                  <div>
                    <p className="text-neutral-80 text-xss">
                      Marketplace Status
                    </p>
                    <p className="font-medium text-neutral-100 text-xs">
                      {data.marketplace_status
                        ?.toLowerCase()
                        .replaceAll("_", " ")
                        .replace(/\b\w/g, (c) => c.toUpperCase())}
                    </p>
                  </div>

                  <div>
                    <p className="text-neutral-80 text-xss">Shipping Status</p>
                    <p className="font-medium text-neutral-100 text-xs">
                      {data.shipping_status
                        ?.toLowerCase()
                        .replaceAll("_", " ")
                        .replace(/\b\w/g, (c) => c.toUpperCase())}
                    </p>
                  </div>

                  <div>
                    <p className="text-neutral-80 text-xss">WMS Status</p>
                    <p className="font-medium text-neutral-100 text-xs">
                      {data.wms_status
                        ?.toLowerCase()
                        .replaceAll("_", " ")
                        .replace(/\b\w/g, (c) => c.toUpperCase())}
                    </p>
                  </div>

                  <div>
                    <p className="text-neutral-80 text-xss">Tracking Number</p>
                    <p className="font-medium text-neutral-100 text-xs">
                      {data.tracking_number ?? "-"}
                    </p>
                  </div>

                  <div>
                    <p className="text-neutral-80 text-xss">Total Amount</p>
                    <p className="font-medium text-neutral-100 text-xs">
                      Rp {data.total_amount.toLocaleString("id-ID")}
                    </p>
                  </div>

                  <div>
                    <p className="text-neutral-80 text-xss">Created At</p>
                    <p className="font-medium text-neutral-100 text-xs">
                      {formatDateTimeSecond(data.created_at)}
                    </p>
                  </div>

                  <div>
                    <p className="text-neutral-80 text-xss">Updated At</p>
                    <p className="font-medium text-neutral-100 text-xs">
                      {formatDateTimeSecond(data.updated_at)}
                    </p>
                  </div>
                </div>

                {/* items table */}
                <div className="overflow-hidden border border-neutral-50 rounded-xl">
                  <table className="min-w-full divide-y text-xs divide-neutral-50">
                    <thead className="bg-white">
                      {table.getHeaderGroups().map((headerGroup) => (
                        <tr key={headerGroup.id}>
                          {headerGroup.headers.map((header) => (
                            <th
                              key={header.id}
                              className={`px-4 py-2 font-medium text-neutral-80 text-left text-xss`}
                            >
                              {flexRender(
                                header.column.columnDef.header,
                                header.getContext(),
                              )}
                            </th>
                          ))}
                        </tr>
                      ))}
                    </thead>

                    <tbody>
                      {table.getRowModel().rows.map((row, i) => (
                        <tr
                          key={row.id}
                          className={`${
                            i % 2 === 0 ? "bg-white" : "bg-neutral-20"
                          }`}
                        >
                          {row.getVisibleCells().map((cell) => (
                            <td key={cell.id} className={`px-4 py-2`}>
                              {flexRender(
                                cell.column.columnDef.cell,
                                cell.getContext(),
                              )}
                            </td>
                          ))}
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>

                {/* action */}
                {currentAction && (
                  <Button
                    variant="default"
                    size="default"
                    className="w-full"
                    onClick={currentAction.action}
                    disabled={currentAction.loading}
                  >
                    {currentAction.loading
                      ? "Processing..."
                      : currentAction.label}
                  </Button>
                )}
              </>
            )}
          </div>
        </>
      )}
    </>
  );
};

export default ActionOrderButton;
