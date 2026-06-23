// ---------- calendar widget ----------
function renderCalendar() {
  const el = document.getElementById('calendar-grid');
  if (!el) return;
  const now = new Date();
  const year = now.getFullYear();
  const month = now.getMonth();
  const today = now.getDate();
  const monthLabel = document.getElementById('calendar-month-label');
  if (monthLabel) {
    monthLabel.textContent = now.toLocaleString('default', { month: 'long' }) + ' ' + year;
  }
  const dows = ['S', 'M', 'T', 'W', 'T', 'F', 'S'];
  let html = dows.map(d => `<div class="dow">${d}</div>`).join('');
  const firstDay = new Date(year, month, 1).getDay();
  const daysInMonth = new Date(year, month + 1, 0).getDate();
  for (let i = 0; i < firstDay; i++) html += `<div class="day empty"></div>`;
  for (let d = 1; d <= daysInMonth; d++) {
    const isToday = d === today;
    html += `<div class="day${isToday ? ' today' : ''}" data-day="${d}" onclick="selectDay(this)">${d}</div>`;
  }
  el.innerHTML = html;
}

function selectDay(node) {
  document.querySelectorAll('#calendar-grid .day.selected').forEach(n => n.classList.remove('selected'));
  node.classList.add('selected');
}

// ---------- pomodoro timer ----------
let pomodoroSeconds = 25 * 60;
let pomodoroInterval = null;

function pomodoroTick() {
  const display = document.getElementById('pomodoro-display');
  if (!display) return;
  const m = Math.floor(pomodoroSeconds / 60).toString().padStart(2, '0');
  const s = (pomodoroSeconds % 60).toString().padStart(2, '0');
  display.textContent = `${m}:${s}`;
}

function pomodoroStart() {
  if (pomodoroInterval) return;
  pomodoroInterval = setInterval(() => {
    if (pomodoroSeconds > 0) {
      pomodoroSeconds--;
      pomodoroTick();
    } else {
      pomodoroPause();
    }
  }, 1000);
}

function pomodoroPause() {
  clearInterval(pomodoroInterval);
  pomodoroInterval = null;
}

function pomodoroReset() {
  pomodoroPause();
  pomodoroSeconds = 25 * 60;
  pomodoroTick();
}

// ---------- GPA calculator ----------
const gpaPoints = { 'A': 4.0, 'A-': 3.7, 'B+': 3.3, 'B': 3.0, 'B-': 2.7, 'C+': 2.3, 'C': 2.0, 'C-': 1.7, 'D': 1.0, 'F': 0.0 };

function gpaAddRow() {
  const tbody = document.getElementById('gpa-rows');
  if (!tbody) return;
  const row = document.createElement('div');
  row.className = 'gpa-row';
  row.innerHTML = `
    <input type="text" placeholder="Course" class="gpa-course">
    <select class="gpa-grade" onchange="gpaCalculate()">
      ${Object.keys(gpaPoints).map(g => `<option value="${g}">${g}</option>`).join('')}
    </select>
    <input type="number" placeholder="Credits" class="gpa-credits" value="3" min="0" step="0.5" oninput="gpaCalculate()">
    <button type="button" class="btn-small btn-danger" onclick="this.parentElement.remove(); gpaCalculate();">✕</button>
  `;
  tbody.appendChild(row);
}

function gpaCalculate() {
  const grades = document.querySelectorAll('.gpa-grade');
  const credits = document.querySelectorAll('.gpa-credits');
  let totalPoints = 0, totalCredits = 0;
  grades.forEach((g, i) => {
    const c = parseFloat(credits[i].value) || 0;
    totalPoints += (gpaPoints[g.value] || 0) * c;
    totalCredits += c;
  });
  const result = document.getElementById('gpa-result');
  if (result) {
    result.textContent = totalCredits > 0 ? 'GPA: ' + (totalPoints / totalCredits).toFixed(2) : 'GPA: --';
  }
}

// ---------- spotify placeholder ----------
function spotifyToggle(btn) {
  const playing = btn.dataset.playing === 'true';
  btn.dataset.playing = (!playing).toString();
  btn.textContent = playing ? '▶' : '⏸';
}

// ---------- generic toggle switch (gym days) ----------
function flipSwitch(el) {
  el.classList.toggle('on');
}

function initGpaIfPresent() {
  const tbody = document.getElementById('gpa-rows');
  if (tbody && tbody.children.length === 0) gpaAddRow();
}

document.addEventListener('DOMContentLoaded', () => {
  renderCalendar();
  pomodoroTick();
  initGpaIfPresent();
});
document.body.addEventListener('htmx:afterSettle', () => {
  renderCalendar();
  initGpaIfPresent();
});
