package expense

import "strconv"

func parseRupiah(value string) int {
	if value == "" {
		return 0
	}

	s := value
	neg := false

	if len(s) > 1 && s[0] == '(' && s[len(s)-1] == ')' {
		neg = true
		s = s[1 : len(s)-1]
	}

	if len(s) > 0 && s[0] == '-' {
		neg = true
		s = s[1:]
	}

	if len(s) > 2 && (s[:2] == "Rp" || s[:2] == "rp" || s[:2] == "RP") {
		s = s[2:]
	}

	clean := ""
	for _, c := range s {
		if c == '.' || c == ',' || c == ' ' {
			continue
		}
		if c < '0' || c > '9' {
			continue
		}
		clean += string(c)
	}

	amount, err := strconv.Atoi(clean)
	if err != nil {
		return 0
	}

	if neg {
		return -amount
	}
	return amount
}

func formatRupiah(amount int) string {
	if amount == 0 {
		return "Rp0"
	}

	neg := amount < 0
	if neg {
		amount = -amount
	}

	s := strconv.Itoa(amount)
	n := len(s)
	if n > 3 {
		result := ""
		for i, c := range s {
			if i > 0 && (n-i)%3 == 0 {
				result += "."
			}
			result += string(c)
		}
		s = result
	}

	if neg {
		return "-Rp" + s
	}
	return "Rp" + s
}

type ExpenseInput struct {
	Year              string `json:"year"`
	Month             string `json:"month"`
	PendapatanUtama   string `json:"pendapatan_utama"`
	PendapatanLainnya string `json:"pendapatan_lainnya"`
	TransferKeNda     string `json:"transfer_ke_nda"`
	TransferKeMimi    string `json:"transfer_ke_mimi"`
	KartuKredit      string `json:"kartu_kredit"`
	TabunganTetap     string `json:"tabungan_tetap"`
	PengeluaranLainnya string `json:"pengeluaran_lainnya"`
}

func (input *ExpenseInput) ToExpense() *Expense {
	year := input.Year
	if year == "" {
		year = strconv.Itoa(2026)
	}

	return &Expense{
		Year:           year,
		Month:          input.Month,
		IncomeMain:     parseRupiah(input.PendapatanUtama),
		IncomeOther:    parseRupiah(input.PendapatanLainnya),
		TransferToNda:  parseRupiah(input.TransferKeNda),
		TransferToMimi: parseRupiah(input.TransferKeMimi),
		CreditCard:     parseRupiah(input.KartuKredit),
		FixedSavings:   parseRupiah(input.TabunganTetap),
		OtherExpense:   parseRupiah(input.PengeluaranLainnya),
	}
}

type ExpenseResponse struct {
	ID            int    `json:"id"`
	Year          string `json:"year"`
	Month         string `json:"month"`
	IncomeMain    string `json:"income_main"`
	IncomeOther   string `json:"income_other"`
	TransferToNda string `json:"transfer_to_nda"`
	TransferToMimi string `json:"transfer_to_mimi"`
	CreditCard    string `json:"credit_card"`
	FixedSavings  string `json:"fixed_savings"`
	OtherExpense  string `json:"other_expense"`
	SisaBulanIni  string `json:"sisa_bulan_ini"`
	TotalSisa     string `json:"total_sisa"`
	TotalTabungan string `json:"total_tabungan"`
	CreatedAt     string `json:"created_at"`
}

func ToExpenseResponse(e *Expense, cumulativeSisa, cumulativeTabungan int) ExpenseResponse {
	income := e.IncomeMain + e.IncomeOther
	expense := e.TransferToNda + e.TransferToMimi + e.CreditCard + e.FixedSavings + e.OtherExpense
	sisa := income - expense
	cumulativeSisa += sisa
	cumulativeTabungan += e.FixedSavings

	return ExpenseResponse{
		ID:            e.ID,
		Year:          e.Year,
		Month:         e.Month,
		IncomeMain:    formatRupiah(e.IncomeMain),
		IncomeOther:   formatRupiah(e.IncomeOther),
		TransferToNda: formatRupiah(e.TransferToNda),
		TransferToMimi: formatRupiah(e.TransferToMimi),
		CreditCard:    formatRupiah(e.CreditCard),
		FixedSavings:  formatRupiah(e.FixedSavings),
		OtherExpense:  formatRupiah(e.OtherExpense),
		SisaBulanIni:  formatRupiah(sisa),
		TotalSisa:     formatRupiah(cumulativeSisa),
		TotalTabungan: formatRupiah(cumulativeTabungan),
		CreatedAt:     e.CreatedAt.Format("2006-01-02T15:04:05"),
	}
}
