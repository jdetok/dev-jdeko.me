// script to load in HTML -- all listener functions are called here 
import * as ui from "./ui.js"
import * as resp from "./resp.js"

export const base = "http://100.81.56.39:8888/bball";

document.addEventListener('DOMContentLoaded', async () => {
    await ui.loadSznOptions();
    await ui.selHvr();
    await ui.randPlayerBtn();
    await ui.search();
    await ui.clearSearch();
    await ui.holdPlayerBtn();
    await resp.getRecentTopScorer();
});