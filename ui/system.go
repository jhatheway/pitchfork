package pitchforkui

import (
	"strconv"
	pf "trident.li/pitchfork/lib"
)

func h_system_settings(cui PfUI) {
	cmd := "system set"
	arg := []string{}

	sys := pf.System_Get()

	msg, err := cui.HandleForm(cmd, arg, sys)

	if msg != "" {
		/* Refresh */
		sys.Refresh()
	}

	var errmsg = ""

	if err != nil {
		/* Failed */
		errmsg = err.Error()
	} else {
		/* Success */
	}

	/* Output the page */
	type Page struct {
		*PfPage
		Message string
		Error   string
		System  pf.PfSys
	}

	p := Page{cui.Page_def(), msg, errmsg, *sys}
	cui.Page_show("system/settings.tmpl", p)
}

func h_system_logA(cui PfUI, user_name string, tg_name string) {
	var err error

	total := 0
	offset := 0
	pageSize := pf.PAGER_PERPAGE /* TODO: Eventually I'd like this to come in from a parameter */

	offset_v, err := cui.FormValue("offset")
	if err == nil && offset_v != "" {
		offset, _ = strconv.Atoi(offset_v)
	}

	search, err := cui.FormValue("search")
	if err != nil {
		search = ""
	}

	var audits []pf.PfAudit
	total, _ = pf.System_AuditMax(search, user_name, tg_name)
	audits, err = pf.System_AuditList(search, user_name, tg_name, offset, pageSize)
	if err != nil {
		cui.Err(err.Error())
		return
	}

	/* Output the page */
	type Page struct {
		*PfPage
		Audits      []pf.PfAudit

		PageSize    int
		LastPage    int
		PagerOffset int
		PagerTotal  int
		Search      string
	}

	/* This math is a bit odd because we depend on the math truncaction. Imagine we're going to the last
	   page of a group of 14 entries, and each page is 4 entries. 14/4 is 3 (truncated). 4*3 gives an offset of 12,
	   the last page of entries 12-14 */
	lastPage := (total/pageSize) * pageSize
	p := Page{cui.Page_def(), audits, pageSize, lastPage, offset, total, search}
	cui.Page_show("system/log.tmpl", p)
}

func h_system_log(cui PfUI) {
	h_system_logA(cui, "", "")
}

func h_system_report(cui PfUI) {
	/* Output the page */
	type Page struct {
		*PfPage
		Message string
	}

	cmd := "system report"
	arg := []string{}
	msg, err := cui.CmdOut(cmd, arg)
	if err != nil {
		msg = err.Error() + "\n" + msg
	}

	p := Page{cui.Page_def(), msg}
	cui.Page_show("system/report.tmpl", p)
}

func h_system_index(cui PfUI) {
	/* Output the page */
	p := cui.Page_def()
	cui.Page_show("system/index.tmpl", p)
}

func h_system(cui PfUI) {
	menu := NewPfUIMenu([]PfUIMentry{
		{"", "", PERM_USER, h_system_index, nil},
		{"log", "Audit Log", PERM_SYS_ADMIN, h_system_log, nil},
		{"report", "Report", PERM_SYS_ADMIN, h_system_report, nil},
		{"settings", "Settings", PERM_SYS_ADMIN, h_system_settings, nil},
		{"iptrk", "IPtrk", PERM_SYS_ADMIN, h_iptrk, nil},
	})

	cui.UIMenu(menu)
}
