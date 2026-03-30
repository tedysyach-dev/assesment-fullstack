import { ArrowDown, ArrowUp } from "lucide-react";

interface StatsCardProps {
  title: string;
  value: number;
  status: boolean;
  statusValue: number;
}

const StatsCard = ({ title, value, status, statusValue }: StatsCardProps) => {
  return (
    <div className="border rounded-lg w-75 border-neutral-50 bg-white flex flex-col gap-2 p-4">
      <h1 className="text-xs text-neutral-80">{title}</h1>

      <div className="flex flex-col">
        <p className="text-bm font-bold">{value}</p>
        <div className="flex flex-row gap-1 items-center">
          {status ? (
            <ArrowUp className="w-4 h-4 text-semantic-green" />
          ) : (
            <ArrowDown className="w-4 h-4 text-semantic-red" />
          )}
          <p
            className={`text-xs font-light ${status ? `text-semantic-green` : `text-semantic-red`}`}
          >
            {statusValue}% this month
          </p>
        </div>
      </div>
    </div>
  );
};

export default StatsCard;
