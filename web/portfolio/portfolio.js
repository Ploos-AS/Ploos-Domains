"use strict";
// UI-only foundation. No login flow, token storage, API calls or demo data.
// A later PR must supply a trusted OIDC session and organization context.
const messages={
 nb:{heading:"Mine domener",intro:"Oversikt over domener som er registrert i porteføljen din.",portfolio:"Domeneportefølje",notConnected:"Ingen tilkoblet konto. Dette er en forhåndsvisning av grensesnittet.",domain:"Domene",registrar:"Registrar",status:"Status",expiry:"Utløpsdato",empty:"Ingen domener å vise.",next:"Neste side"},
 en:{heading:"My domains",intro:"Domains tracked in your portfolio.",portfolio:"Domain portfolio",notConnected:"No connected account. This is a preview of the interface.",domain:"Domain",registrar:"Registrar",status:"Status",expiry:"Expiry date",empty:"No domains to display.",next:"Next page"}
};
let lang="nb";
function translate(){
 document.documentElement.lang=lang;
 document.querySelectorAll("[data-i18n]").forEach(el=>{el.textContent=messages[lang][el.dataset.i18n]});
 const button=document.getElementById("language");
 button.textContent=lang==="nb"?"English":"Norsk";
 button.setAttribute("aria-label",lang==="nb"?"Switch to English":"Bytt til norsk");
}
document.getElementById("language").addEventListener("click",()=>{lang=lang==="nb"?"en":"nb";translate()});
translate();
