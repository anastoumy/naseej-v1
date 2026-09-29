import { useState } from "react";
import {
  Calculator,
  TrendingUp,
  CreditCard,
  Banknote,
  AlertTriangle,
  CheckCircle,
  Printer,
  Lock,
} from "lucide-react";
import { useTranslation } from "react-i18next";
import { DUMMY_DAILY_SUMMARY, MERCHANT_CONFIG } from "@/data/dummy-data";
import { cn, formatSAR, formatDate } from "@/lib/utils";

export default function EndOfDayPage() {
  const { t } = useTranslation();
  const summary = DUMMY_DAILY_SUMMARY;
  const [countedCash, setCountedCash] = useState("");
  const [safeDropAmount, setSafeDropAmount] = useState("");
  const [reportGenerated, setReportGenerated] = useState(false);

  const expectedCash =
    summary.openingCash + summary.cashSales - summary.refundsTotal;
  const counted = parseFloat(countedCash) || 0;
  const variance = counted - expectedCash;
  const hasVariance = Math.abs(variance) > 0.01;

  const handleGenerateReport = () => {
    setReportGenerated(true);
  };

  return (
    <div className="flex-1 flex flex-col overflow-hidden">
      {/* Header */}
      <div className="h-12 bg-card border-b border-border flex items-center gap-3 px-4 shrink-0">
        <Calculator className="w-4 h-4 text-muted-foreground" />
        <h1 className="text-sm font-semibold text-foreground">
          {t('features.pos.screens.endOfDay.title')}
        </h1>
        <span className="badge-muted ml-2">{formatDate(summary.date)}</span>
        <span className="badge bg-accent/15 text-accent ml-1">
          {MERCHANT_CONFIG.location}
        </span>
      </div>

      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {/* Z-Report Summary */}
        <div>
          <h2 className="text-xs font-semibold text-muted-foreground uppercase tracking-widest mb-3">
            {t('features.pos.screens.endOfDay.zReportSummary')}
          </h2>
          <div className="grid grid-cols-2 gap-3">
            {/* Total sales */}
            <div className="pos-card p-4">
              <div className="flex items-center gap-2 mb-2">
                <TrendingUp className="w-4 h-4 text-success" />
                <span className="text-xs text-muted-foreground">
                  {t('features.pos.screens.endOfDay.totalSales')}
                </span>
              </div>
              <p className="text-2xl font-bold text-success tabular-nums">
                {formatSAR(summary.totalSales)}
              </p>
            </div>

            {/* Transactions */}
            <div className="pos-card p-4">
              <div className="flex items-center gap-2 mb-2">
                <Calculator className="w-4 h-4 text-accent-blue" />
                <span className="text-xs text-muted-foreground">
                  Transactions
                </span>
              </div>
              <p className="text-2xl font-bold text-foreground tabular-nums">
                {summary.transactionCount}
              </p>
              <p className="text-xs text-muted-foreground mt-0.5">
                {summary.itemsSold} {t('features.pos.screens.endOfDay.itemsSold')}
              </p>
            </div>

            {/* Cash sales */}
            <div className="pos-card p-4">
              <div className="flex items-center gap-2 mb-2">
                <Banknote className="w-4 h-4 text-success" />
                <span className="text-xs text-muted-foreground">
                  {t('features.pos.screens.endOfDay.cashSales')}
                </span>
              </div>
              <p className="text-xl font-bold text-foreground tabular-nums">
                {formatSAR(summary.cashSales)}
              </p>
            </div>

            {/* Card sales */}
            <div className="pos-card p-4">
              <div className="flex items-center gap-2 mb-2">
                <CreditCard className="w-4 h-4 text-purple-400" />
                <span className="text-xs text-muted-foreground">
                  {t('features.pos.screens.endOfDay.cardSales')}
                </span>
              </div>
              <p className="text-xl font-bold text-foreground tabular-nums">
                {formatSAR(summary.cardSales)}
              </p>
            </div>
          </div>

          {/* Refunds row */}
          <div className="pos-card p-4 mt-3 flex items-center justify-between">
            <div>
              <p className="text-xs text-muted-foreground">{t('features.pos.screens.endOfDay.totalRefunds')}</p>
              <p className="text-lg font-semibold text-destructive tabular-nums mt-0.5">
                -{formatSAR(summary.refundsTotal)}
              </p>
            </div>
            <div className="text-right">
              <p className="text-xs text-muted-foreground">{t('features.pos.screens.endOfDay.netRevenue')}</p>
              <p className="text-lg font-semibold text-success tabular-nums mt-0.5">
                {formatSAR(summary.totalSales - summary.refundsTotal)}
              </p>
            </div>
          </div>
        </div>

        {/* Cash Reconciliation */}
        <div>
          <h2 className="text-xs font-semibold text-muted-foreground uppercase tracking-widest mb-3">
            {t('features.pos.screens.endOfDay.cashReconciliation')}
          </h2>
          <div className="pos-card p-4 space-y-3">
            <div className="grid grid-cols-2 gap-3 text-sm">
              <div className="flex justify-between">
                <span className="text-muted-foreground">{t('features.pos.screens.endOfDay.openingCash')}</span>
                <span className="text-foreground tabular-nums">
                  {formatSAR(summary.openingCash)}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">{t('features.pos.screens.endOfDay.plusCashSales')}</span>
                <span className="text-success tabular-nums">
                  +{formatSAR(summary.cashSales)}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">{t('features.pos.screens.endOfDay.minusCashRefunds')}</span>
                <span className="text-destructive tabular-nums">
                  -{formatSAR(summary.refundsTotal)}
                </span>
              </div>
              <div className="flex justify-between font-semibold">
                <span className="text-muted-foreground">{t('features.pos.screens.endOfDay.expected')}</span>
                <span className="text-foreground tabular-nums">
                  {formatSAR(expectedCash)}
                </span>
              </div>
            </div>

            <div className="h-px bg-border" />

            {/* Count input */}
            <div>
              <label className="block text-sm font-medium text-foreground mb-2">
                {t('features.pos.screens.endOfDay.actualCashCount')}
              </label>
              <input
                type="number"
                value={countedCash}
                onChange={(e) => setCountedCash(e.target.value)}
                placeholder={`Enter counted amount (Expected: ${expectedCash.toFixed(2)})`}
                className="pos-input"
                min={0}
                step={0.01}
                aria-label="Actual cash counted"
              />
            </div>

            {/* Variance display */}
            {countedCash && (
              <div
                className={cn(
                  "rounded-lg p-3 flex items-center justify-between",
                  !hasVariance
                    ? "bg-success/10 border border-success/20"
                    : variance > 0
                      ? "bg-warning/10 border border-warning/20"
                      : "bg-destructive/10 border border-destructive/20",
                )}
              >
                <div className="flex items-center gap-2">
                  {!hasVariance ? (
                    <CheckCircle className="w-4 h-4 text-success" />
                  ) : (
                    <AlertTriangle className="w-4 h-4 text-warning" />
                  )}
                  <span className="text-sm font-medium text-foreground">
                    {!hasVariance
                      ? t('features.pos.screens.endOfDay.balanced')
                      : variance > 0
                        ? t('features.pos.screens.endOfDay.cashOver')
                        : t('features.pos.screens.endOfDay.cashShort')}
                  </span>
                </div>
                <span
                  className={cn(
                    "text-lg font-bold tabular-nums",
                    !hasVariance
                      ? "text-success"
                      : variance > 0
                        ? "text-warning"
                        : "text-destructive",
                  )}
                >
                  {variance >= 0 ? "+" : ""}
                  {formatSAR(variance)}
                </span>
              </div>
            )}
          </div>
        </div>

        {/* Safe Drop */}
        <div>
          <h2 className="text-xs font-semibold text-muted-foreground uppercase tracking-widest mb-3">
            Safe Drop
          </h2>
          <div className="pos-card p-4">
            <label className="block text-sm font-medium text-foreground mb-2">
              Cash Removed for Deposit
            </label>
            <input
              type="number"
              value={safeDropAmount}
              onChange={(e) => setSafeDropAmount(e.target.value)}
              placeholder="Enter safe drop amount"
              className="pos-input"
              min={0}
              step={0.01}
              aria-label="Safe drop amount"
            />
            {safeDropAmount && (
              <p className="text-xs text-muted-foreground mt-2">
                Remaining in drawer:{" "}
                <span className="text-foreground font-medium tabular-nums">
                  {formatSAR(Math.max(0, counted - parseFloat(safeDropAmount)))}
                </span>
              </p>
            )}
          </div>
        </div>

        {/* Actions */}
        <div className="flex gap-3 pb-4">
          <button
            onClick={handleGenerateReport}
            className={cn(
              "flex-1 flex items-center justify-center gap-2 h-12 rounded-xl text-sm font-semibold transition-all duration-150 cursor-pointer",
              "focus:outline-none focus:ring-2 focus:ring-accent/50",
              reportGenerated
                ? "bg-success/10 border border-success/20 text-success"
                : "bg-secondary text-foreground hover:bg-secondary/80",
            )}
          >
            <Printer className="w-4 h-4" />
            {reportGenerated ? "Z-Report Printed" : "Print Z-Report"}
          </button>

          <button
            className="flex-1 flex items-center justify-center gap-2 h-12 rounded-xl text-sm font-semibold bg-destructive text-white hover:bg-destructive/90 active:scale-[0.98] transition-all duration-150 cursor-pointer focus:outline-none focus:ring-2 focus:ring-destructive/50"
            onClick={() => alert("Manager PIN required to close shift")}
          >
            <Lock className="w-4 h-4" />
            Close Shift
          </button>
        </div>
      </div>
    </div>
  );
}
