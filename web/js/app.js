/**
 * Narxoz College (NC) — веб-прототип
 * Палитра: #a82523, #d50032 (сайт университета)
 */

const API_BASE = (typeof window !== 'undefined' && window.location?.origin && window.location.origin.startsWith('http'))
  ? window.location.origin
  : 'http://localhost:8080';

let state = {
  token: null,
  user: null,
  currentPage: 'home',
  chatWs: null,
  chatMessages: [],
  forumTopicId: null,
  forumTopic: null,
  notesNoteId: null,
  notesNote: null,
  directorPostId: null,
  directorPost: null
};

// ——— API ———
function getAuthHeaders(includeContentType = true) {
  const h = {};
  if (includeContentType) h['Content-Type'] = 'application/json';
  if (state.token) h['Authorization'] = 'Bearer ' + state.token;
  return h;
}

async function api(path, options = {}) {
  const url = path.startsWith('http') ? path : API_BASE + path;
  const hasBody = options.body != null && options.body !== '';
  const res = await fetch(url, {
    ...options,
    headers: { ...getAuthHeaders(hasBody), ...(options.headers || {}) }
  });
  if (res.status === 401) {
    logout();
    window.location.reload();
    throw new Error('Unauthorized');
  }
  const text = await res.text();
  let data = null;
  if (text) {
    try { data = JSON.parse(text); } catch (_) {}
  }
  if (!res.ok) {
    throw new Error(data?.error || res.statusText || 'Ошибка');
  }
  return data;
}

// ——— Хранилище ———
function loadStored() {
  try {
    const t = localStorage.getItem('nc_token');
    const u = localStorage.getItem('nc_user');
    if (t && u) {
      state.token = t;
      state.user = JSON.parse(u);
      return true;
    }
  } catch (_) {}
  return false;
}

function saveStored() {
  if (state.token) localStorage.setItem('nc_token', state.token);
  if (state.user) localStorage.setItem('nc_user', JSON.stringify(state.user));
}

function clearStored() {
  localStorage.removeItem('nc_token');
  localStorage.removeItem('nc_user');
}

// ——— Навигация ———
const PAGES = [
  { id: 'home', title: 'Главная', icon: '', roles: ['student', 'teacher', 'director', 'admin'] },
  { id: 'schedule', title: 'Расписание', icon: '', roles: ['student', 'teacher'] },
  { id: 'grades', title: 'Оценки', icon: '', roles: ['student'] },
  { id: 'groups', title: 'Группы', icon: '', roles: ['teacher', 'director', 'admin'] },
  { id: 'teacherGrades', title: 'Выставить оценку', icon: '', roles: ['teacher'] },
  { id: 'teacherHomework', title: 'ДЗ к парам', icon: '', roles: ['teacher'] },
  { id: 'forum', title: 'Форум', icon: '', roles: ['student', 'teacher', 'director', 'admin'] },
  { id: 'events', title: 'События', icon: '', roles: ['student', 'teacher', 'director', 'admin'] },
  { id: 'notes', title: 'Мои заметки', icon: '', roles: ['student', 'teacher', 'director', 'admin'] },
  { id: 'directorProfile', title: 'Профиль директора', icon: '', roles: ['student', 'teacher', 'director'] },
  { id: 'regulations', title: 'Регламент', icon: '', roles: ['student', 'teacher', 'director', 'admin'] },
  { id: 'map', title: 'Карта здания', icon: '', roles: ['student', 'teacher', 'director', 'admin'] },
  { id: 'profile', title: 'Профиль', icon: '', roles: ['student', 'teacher', 'director', 'admin'] }
];

function navItems() {
  const role = state.user?.role || '';
  return PAGES.filter(p => p.roles.includes(role));
}

function renderNav() {
  const items = navItems();
  const nav = document.getElementById('nav');
  if (nav) {
    nav.innerHTML = items.map(p => `
      <button class="nav-item ${state.currentPage === p.id ? 'active' : ''}" data-page="${p.id}">
        ${p.icon ? `<span>${p.icon}</span>` : ''}
        <span style="${p.icon ? 'margin-left: 0.5rem' : ''}">${p.title}</span>
      </button>
    `).join('');
    nav.querySelectorAll('.nav-item').forEach(el => {
      el.addEventListener('click', () => goTo(el.dataset.page));
    });
  }
  const bottomInner = document.getElementById('bottom-nav-inner');
  if (bottomInner) {
    bottomInner.innerHTML = items.map(p => `
      <button type="button" class="nav-item ${state.currentPage === p.id ? 'active' : ''}" data-page="${p.id}">
        ${p.icon ? `<span>${p.icon}</span>` : ''}
        <span>${p.title}</span>
      </button>
    `).join('');
    bottomInner.querySelectorAll('.nav-item').forEach(el => {
      el.addEventListener('click', () => {
        goTo(el.dataset.page);
        document.getElementById('sidebar').classList.remove('open');
        var ov = document.getElementById('sidebar-overlay');
        if (ov) { ov.classList.add('hidden'); ov.classList.remove('visible'); }
      });
    });
  }
}

function goTo(pageId) {
  state.currentPage = pageId;
  const page = PAGES.find(p => p.id === pageId);
  document.getElementById('page-title').textContent = page ? page.title : pageId;
  renderNav();
  const content = document.getElementById('content');
  content.innerHTML = '<p class="empty-state">Загрузка…</p>';
  if (pageId === 'home') loadHome();
  else if (pageId === 'schedule') loadSchedule();
  else if (pageId === 'grades') loadGrades();
  else if (pageId === 'groups') loadGroups();
  else if (pageId === 'teacherGrades') loadTeacherGrades();
  else if (pageId === 'teacherHomework') loadTeacherHomework();
  else if (pageId === 'forum') loadForum();
  else if (pageId === 'events') loadEvents();
  else if (pageId === 'notes') loadNotes();
  else if (pageId === 'directorProfile') loadDirectorProfile();
  else if (pageId === 'regulations') loadRegulations();
  else if (pageId === 'map') loadMap();
  else if (pageId === 'profile') loadProfile();
}

// ——— Экраны ———
function loadHome() {
  api('/api/events?limit=10').then(list => {
    document.getElementById('content').innerHTML = `
      <div class="card">
        <h3>Новости и мероприятия</h3>
        ${Array.isArray(list) && list.length
          ? list.map(e => `
            <div class="list-item">
              <div>
                <strong>${escapeHtml(e.title)}</strong>
                ${e.description ? `<p class="meta">${escapeHtml(e.description)}</p>` : ''}
                <p class="meta">${formatDate(e.event_date)} ${e.location ? ' · ' + escapeHtml(e.location) : ''}</p>
              </div>
            </div>
          `).join('')
          : '<p class="empty-state">Нет событий</p>'}
      </div>
    `;
  }).catch(err => {
    document.getElementById('content').innerHTML = `<p class="error-msg">${escapeHtml(err.message)}</p>`;
  });
}

function loadSchedule() {
  api('/api/schedule').then(list => {
    const byDay = {};
    (list || []).forEach(s => {
      const d = s.day_of_week || 0;
      if (!byDay[d]) byDay[d] = [];
      byDay[d].push(s);
    });
    const days = ['', 'Понедельник', 'Вторник', 'Среда', 'Четверг', 'Пятница', 'Суббота', 'Воскресенье'];
    document.getElementById('content').innerHTML = `
      <div class="card">
        <h3>Расписание</h3>
        ${[1,2,3,4,5,6,7].map(d => {
          const items = byDay[d] || [];
          if (!items.length) return '';
          return `
            <div style="margin-bottom: 1.25rem">
              <strong>${days[d]}</strong>
              ${items.map(s => `
                <div class="list-item schedule-item">
                  <div style="flex:1">
                    <strong>${escapeHtml(s.subject)}</strong>
                    <p class="meta">${formatTimeStr(s.start_time)} – ${formatTimeStr(s.end_time)} · Ауд. ${escapeHtml(s.room || '—')}</p>
                    <button type="button" class="btn btn-text btn-small btn-attachments" data-schedule-id="${s.id}">ДЗ и материалы</button>
                    <div id="attachments-${s.id}" class="attachments-list hidden"></div>
                  </div>
                </div>
              `).join('')}
            </div>
          `;
        }).filter(Boolean).join('') || '<p class="empty-state">Нет занятий</p>'}
      </div>
    `;
    document.querySelectorAll('.btn-attachments').forEach(btn => {
      btn.addEventListener('click', function () {
        const id = this.dataset.scheduleId;
        const wrap = document.getElementById('attachments-' + id);
        if (wrap.classList.contains('loaded')) {
          wrap.classList.toggle('hidden');
          return;
        }
        api('/api/schedule/' + id + '/attachments').then(att => {
          wrap.innerHTML = (att && att.length) ? att.map(a => `
            <div class="attachment-item"><strong>${escapeHtml(a.title)}</strong> (${a.attachment_type || 'homework'})<br><span class="meta">${escapeHtml((a.body || '').slice(0, 200))}</span></div>
          `).join('') : '<p class="meta">Нет материалов</p>';
          wrap.classList.add('loaded');
          wrap.classList.remove('hidden');
        }).catch(() => { wrap.innerHTML = '<p class="meta">Ошибка загрузки</p>'; wrap.classList.remove('hidden'); });
      });
    });
  }).catch(err => {
    document.getElementById('content').innerHTML = `<p class="error-msg">${escapeHtml(err.message)}</p>`;
  });
}
function formatTimeStr(t) {
  if (!t) return '';
  if (typeof t === 'string') return t.slice(0, 5);
  return String(t);
}

function loadGrades() {
  Promise.all([api('/api/grades'), api('/api/grades/gpa')]).then(([list, gpaRes]) => {
    const gpa = gpaRes?.gpa != null ? gpaRes.gpa : '-';
    document.getElementById('content').innerHTML = `
      <div class="card">
        <h3>Журнал оценок</h3>
        <p><strong>Средний балл:</strong> ${gpa}</p>
        <p class="meta">Оценки по 100-балльной шкале</p>
        ${Array.isArray(list) && list.length
          ? list.map(g => `
            <div class="list-item">
              <div>
                <strong>${escapeHtml(g.subject)}</strong>
                <p class="meta">${formatDate(g.grade_date)} — ${g.grade} баллов ${g.comment ? ' · ' + escapeHtml(g.comment) : ''}</p>
              </div>
            </div>
          `).join('')
          : '<p class="empty-state">Нет оценок</p>'}
      </div>
    `;
  }).catch(err => {
    document.getElementById('content').innerHTML = `<p class="error-msg">${escapeHtml(err.message)}</p>`;
  });
}

function loadTeacherGrades() {
  Promise.all([api('/api/groups'), api('/api/schedule')]).then(([groups, schedule]) => {
    const gr = Array.isArray(groups) ? groups : [];
    const sch = Array.isArray(schedule) ? schedule : [];
    document.getElementById('content').innerHTML = `
      <div class="card">
        <h3>Выставить оценку</h3>
        <div class="form-row">
          <label>Группа</label>
          <select id="tg-group"><option value="">— Выберите группу —</option>${gr.map(g => `<option value="${g.id}">${escapeHtml(g.name)}</option>`).join('')}</select>
        </div>
        <div class="form-row">
          <label>Студент</label>
          <select id="tg-student"><option value="">— Сначала выберите группу —</option></select>
        </div>
        <div class="form-row">
          <label>Занятие (предмет)</label>
          <select id="tg-schedule"><option value="">— Выберите —</option>${sch.map(s => `<option value="${s.id}">${escapeHtml(s.subject)} · ${escapeHtml(s.room || '')}</option>`).join('')}</select>
        </div>
        <div class="form-row">
          <label>Оценка</label>
          <input type="number" id="tg-grade" min="0" max="100" step="0.5" placeholder="85">
        </div>
        <div class="form-row">
          <label>Дата (ДД/ММ/ГГГГ)</label>
          <input type="text" id="tg-date" placeholder="31/01/2026" value="${formatDate(new Date().toISOString())}">
        </div>
        <div class="form-row">
          <label>Комментарий</label>
          <input type="text" id="tg-comment" placeholder="Необязательно">
        </div>
        <button class="btn btn-primary" id="tg-submit">Выставить оценку</button>
      </div>
    `;
    const groupSel = document.getElementById('tg-group');
    const studentSel = document.getElementById('tg-student');
    groupSel.addEventListener('change', () => {
      const gid = groupSel.value;
      studentSel.innerHTML = '<option value="">— Выберите студента —</option>';
      if (!gid) return;
      api('/api/groups/' + gid + '/students').then(students => {
        (students || []).forEach(s => {
          studentSel.innerHTML += `<option value="${s.user_id}">${escapeHtml(s.full_name)}</option>`;
        });
      });
    });
    document.getElementById('tg-submit').addEventListener('click', () => {
      const userId = studentSel.value;
      const scheduleId = document.getElementById('tg-schedule').value;
      const grade = parseFloat(document.getElementById('tg-grade').value);
      const gradeDateInput = document.getElementById('tg-date').value.trim();
      const gradeDate = parseDDMMYYYY(gradeDateInput) || gradeDateInput;
      const comment = document.getElementById('tg-comment').value.trim();
      if (!userId || !scheduleId || isNaN(grade)) { alert('Заполните группу, студента, занятие и оценку'); return; }
      if (!gradeDate) { alert('Введите дату в формате ДД/ММ/ГГГГ'); return; }
      api('/api/grades', {
        method: 'POST',
        body: JSON.stringify({ user_id: parseInt(userId, 10), schedule_id: parseInt(scheduleId, 10), grade, grade_date: gradeDate, comment: comment || '' })
      }).then(() => { alert('Оценка выставлена'); loadTeacherGrades(); }).catch(err => alert(err.message));
    });
  }).catch(err => {
    document.getElementById('content').innerHTML = `<p class="error-msg">${escapeHtml(err.message)}</p>`;
  });
}

function loadTeacherHomework() {
  api('/api/schedule').then(list => {
    const sch = Array.isArray(list) ? list : [];
    document.getElementById('content').innerHTML = `
      <div class="card">
        <h3>ДЗ и материалы к парам</h3>
        <p class="meta">Выберите занятие и добавьте ДЗ/уведомление или просмотрите существующие.</p>
        ${sch.length ? sch.map(s => `
          <div class="list-item" style="flex-wrap:wrap">
            <div style="flex:1; min-width:200px">
              <strong>${escapeHtml(s.subject)}</strong>
              <p class="meta">${formatTimeStr(s.start_time)} · Ауд. ${escapeHtml(s.room || '—')} · ID ${s.id}</p>
              <div id="hw-list-${s.id}" class="attachments-list"></div>
              <button type="button" class="btn btn-primary btn-small btn-add-hw" data-schedule-id="${s.id}">+ Добавить ДЗ</button>
              <div id="hw-form-${s.id}" class="card hidden" style="margin-top:0.5rem">
                <input type="text" id="hw-title-${s.id}" placeholder="Название">
                <textarea id="hw-body-${s.id}" placeholder="Текст" rows="2"></textarea>
                <select id="hw-type-${s.id}"><option value="homework">ДЗ</option><option value="notification">Уведомление</option><option value="material">Материал</option></select>
                <button type="button" class="btn btn-small btn-primary btn-save-hw" data-schedule-id="${s.id}">Сохранить</button>
                <button type="button" class="btn btn-small btn-text btn-cancel-hw" data-schedule-id="${s.id}">Отмена</button>
              </div>
            </div>
          </div>
        `).join('') : '<p class="empty-state">Нет занятий в расписании</p>'}
      </div>
    `;
    function loadAttachments(scheduleId) {
      api('/api/schedule/' + scheduleId + '/attachments').then(att => {
        const wrap = document.getElementById('hw-list-' + scheduleId);
        if (!wrap) return;
        wrap.innerHTML = (att && att.length) ? att.map(a => `
          <div class="attachment-item" style="display:flex;justify-content:space-between;align-items:start;margin:0.5rem 0">
            <div><strong>${escapeHtml(a.title)}</strong> (${a.attachment_type || 'homework'})<br><span class="meta">${escapeHtml((a.body || '').slice(0, 100))}</span></div>
            <div>
              <button type="button" class="btn btn-small btn-danger btn-del-att" data-att-id="${a.id}" data-schedule-id="${scheduleId}">Удалить</button>
            </div>
          </div>
        `).join('') : '';
        wrap.querySelectorAll('.btn-del-att').forEach(btn => {
          btn.addEventListener('click', () => {
            if (!confirm('Удалить?')) return;
            api('/api/attachments/' + btn.dataset.attId, { method: 'DELETE' }).then(() => loadAttachments(btn.dataset.scheduleId)).catch(err => alert(err.message));
          });
        });
      });
    }
    sch.forEach(s => { loadAttachments(s.id); });
    document.querySelectorAll('.btn-add-hw').forEach(btn => {
      btn.addEventListener('click', () => {
        const id = btn.dataset.scheduleId;
        document.getElementById('hw-form-' + id).classList.toggle('hidden');
      });
    });
    document.querySelectorAll('.btn-cancel-hw').forEach(btn => {
      btn.addEventListener('click', () => document.getElementById('hw-form-' + btn.dataset.scheduleId).classList.add('hidden'));
    });
    document.querySelectorAll('.btn-save-hw').forEach(btn => {
      btn.addEventListener('click', () => {
        const sid = btn.dataset.scheduleId;
        const title = document.getElementById('hw-title-' + sid).value.trim();
        const body = document.getElementById('hw-body-' + sid).value.trim();
        const type = document.getElementById('hw-type-' + sid).value || 'homework';
        if (!title) { alert('Введите название'); return; }
        api('/api/schedule/' + sid + '/attachments', {
          method: 'POST',
          body: JSON.stringify({ title, body, attachment_type: type })
        }).then(() => {
          document.getElementById('hw-form-' + sid).classList.add('hidden');
          document.getElementById('hw-title-' + sid).value = '';
          document.getElementById('hw-body-' + sid).value = '';
          loadAttachments(sid);
        }).catch(err => alert(err.message));
      });
    });
  }).catch(err => {
    document.getElementById('content').innerHTML = `<p class="error-msg">${escapeHtml(err.message)}</p>`;
  });
}

function loadGroups() {
  api('/api/groups').then(list => {
    document.getElementById('content').innerHTML = `
      <div class="card">
        <h3>Группы</h3>
        ${Array.isArray(list) && list.length
          ? list.map(g => `
            <div class="list-item">
              <div>
                <strong>${escapeHtml(g.name)}</strong>
                <p class="meta">${escapeHtml(g.description || '')} · Студентов: ${g.student_count ?? 0}</p>
              </div>
              <button class="btn btn-small btn-secondary btn-group-students" data-group-id="${g.id}" data-group-name="${escapeHtml(g.name)}">Студенты</button>
            </div>
          `).join('')
          : '<p class="empty-state">Нет групп</p>'}
      </div>
      <div id="students-panel" class="card hidden">
        <h3 id="students-panel-title">Студенты группы</h3>
        <div id="students-panel-list"></div>
        <button type="button" class="btn btn-text btn-small" id="students-panel-close">Закрыть</button>
      </div>
    `;
    document.getElementById('content').querySelectorAll('.btn-group-students').forEach(btn => {
      btn.addEventListener('click', () => {
        const id = btn.dataset.groupId;
        const name = btn.dataset.groupName || '';
        api('/api/groups/' + id + '/students').then(students => {
          const panel = document.getElementById('students-panel');
          const listEl = document.getElementById('students-panel-list');
          document.getElementById('students-panel-title').textContent = 'Студенты: ' + name;
          listEl.innerHTML = Array.isArray(students) && students.length
            ? students.map(s => `<div class="list-item"><div><strong>${escapeHtml(s.full_name)}</strong><p class="meta">${escapeHtml(s.email || '')}</p></div></div>`).join('')
            : '<p class="empty-state">Нет студентов</p>';
          panel.classList.remove('hidden');
        }).catch(err => {
          alert(err.message);
        });
      });
    });
    document.getElementById('students-panel-close').addEventListener('click', () => {
      document.getElementById('students-panel').classList.add('hidden');
    });
  }).catch(err => {
    document.getElementById('content').innerHTML = `<p class="error-msg">${escapeHtml(err.message)}</p>`;
  });
}

function loadLibrary() {
  api('/api/library/books?limit=20').then(list => {
    document.getElementById('content').innerHTML = `
      <div class="form-row">
        <input type="text" id="lib-search" placeholder="Поиск по названию или автору..." style="max-width: 320px">
        <button class="btn btn-primary btn-small" style="margin-left: 0.5rem" id="lib-search-btn">Искать</button>
      </div>
      <div class="card" id="library-list">
        <h3>Книги</h3>
        ${Array.isArray(list) && list.length
          ? list.map(b => `
            <div class="list-item">
              <div>
                <strong>${escapeHtml(b.title)}</strong>
                <p class="meta">${escapeHtml(b.author || '')} · Доступно: ${b.available ?? 0} из ${b.total_copies ?? 0}</p>
              </div>
              <button class="btn btn-small btn-primary" data-book-id="${b.id}">Забронировать</button>
            </div>
          `).join('')
          : '<p class="empty-state">Нет книг</p>'}
      </div>
    `;
    document.getElementById('lib-search-btn').addEventListener('click', () => {
      const q = document.getElementById('lib-search').value.trim();
      api('/api/library/books?q=' + encodeURIComponent(q) + '&limit=20').then(books => {
        const wrap = document.getElementById('library-list');
        wrap.innerHTML = '<h3>Книги</h3>' + (Array.isArray(books) && books.length
          ? books.map(b => `
            <div class="list-item">
              <div>
                <strong>${escapeHtml(b.title)}</strong>
                <p class="meta">${escapeHtml(b.author || '')} · Доступно: ${b.available ?? 0}</p>
              </div>
              <button class="btn btn-small btn-primary" data-book-id="${b.id}">Забронировать</button>
            </div>
          `).join('')
          : '<p class="empty-state">Ничего не найдено</p>');
        wrap.querySelectorAll('[data-book-id]').forEach(btn => {
          btn.addEventListener('click', function() { reserveBook(this.dataset.bookId, this); });
        });
      });
    });
    document.getElementById('content').querySelectorAll('[data-book-id]').forEach(btn => {
      btn.addEventListener('click', function() { reserveBook(this.dataset.bookId, this); });
    });
  }).catch(err => {
    document.getElementById('content').innerHTML = `<p class="error-msg">${escapeHtml(err.message)}</p>`;
  });
}

function reserveBook(bookId, btnEl) {
  if (btnEl) {
    btnEl.disabled = true;
    btnEl.textContent = '...';
  }
  api('/api/library/reservations', {
    method: 'POST',
    body: JSON.stringify({ book_id: parseInt(bookId, 10) })
  }).then(() => {
    alert('Бронь создана');
    loadLibrary();
  }).catch(err => {
    alert(err.message || 'Ошибка бронирования');
    if (btnEl) {
      btnEl.disabled = false;
      btnEl.textContent = 'Забронировать';
    }
  });
}

function loadForum() {
  if (state.forumTopicId) {
    const topic = state.forumTopic || {};
    api('/api/forum/posts?parent_id=' + state.forumTopicId + '&limit=50').then(replies => {
      const reps = Array.isArray(replies) ? replies : [];
      document.getElementById('content').innerHTML = `
        <div class="card">
          <button type="button" class="btn btn-text btn-small" id="forum-back">← Назад к темам</button>
          <h3>${escapeHtml(topic.title || 'Тема')}</h3>
          <div class="forum-topic-body">${escapeHtml(topic.body || '')}</div>
          <p class="meta">${topic.is_anonymous ? 'Анонимно' : (topic.author_name || '')} · ${formatDate(topic.created_at)}</p>
          <h4 style="margin-top:1rem">Ответы (${reps.length})</h4>
          <div id="forum-replies">${reps.map(r => `
            <div class="list-item forum-reply-item" data-post-id="${r.id}">
              <div>${r.media_url ? `<img src="${API_BASE}${r.media_url}" alt="" style="max-width:200px;max-height:150px;border-radius:2px">` : ''}${escapeHtml(r.body || '')}</div>
              <p class="meta">${r.is_anonymous ? 'Анонимно' : (r.author_id ? `<span class="forum-author-click" data-user-id="${r.author_id}">${escapeHtml(r.author_name || '')}</span>` : escapeHtml(r.author_name || ''))} · ${formatDate(r.created_at)}</p>
            </div>
          `).join('')}</div>
          <div class="card" style="margin-top:1rem">
            <h4>Ответить</h4>
            <div class="form-row"><label>Текст</label><textarea id="forum-reply-body" placeholder="Текст ответа"></textarea></div>
            <div class="form-row"><label><input type="checkbox" id="forum-reply-anonymous"> Публиковать анонимно</label></div>
            <button type="button" class="btn btn-primary" id="forum-reply-submit">Отправить</button>
          </div>
        </div>
      `;
      document.getElementById('forum-back').addEventListener('click', () => {
        state.forumTopicId = null;
        loadForum();
      });
      document.getElementById('forum-reply-submit').addEventListener('click', () => {
        const body = document.getElementById('forum-reply-body').value.trim();
        const isAnonymous = document.getElementById('forum-reply-anonymous').checked;
        if (!body) { alert('Введите текст'); return; }
        api('/api/forum/posts', {
          method: 'POST',
          body: JSON.stringify({ parent_id: state.forumTopicId, title: null, body, is_anonymous: isAnonymous })
        }).then(() => { document.getElementById('forum-reply-body').value = ''; loadForum(); }).catch(err => alert(err.message));
      });
      document.querySelectorAll('.forum-author-click').forEach(el => {
        el.style.cursor = 'pointer';
        el.addEventListener('click', () => showUserProfile(el.dataset.userId));
      });
    }).catch(err => {
      document.getElementById('content').innerHTML = `<p class="error-msg">${escapeHtml(err.message)}</p>`;
    });
    return;
  }
  state.forumTopic = null;
  api('/api/forum/posts?limit=30').then(list => {
    document.getElementById('content').innerHTML = `
      <div class="card">
        <h3>Темы форума</h3>
        <button class="btn btn-primary btn-small" id="forum-new-topic">Новая тема</button>
        <div id="forum-topics" style="margin-top: 1rem"></div>
      </div>
      <div id="forum-form" class="card hidden">
        <h3>Новая тема</h3>
        <div class="form-row">
          <label>Заголовок</label>
          <input type="text" id="forum-title" placeholder="Заголовок">
        </div>
        <div class="form-row">
          <label>Текст</label>
          <textarea id="forum-body" placeholder="Текст"></textarea>
        </div>
        <div class="form-row">
          <label><input type="checkbox" id="forum-anonymous"> Публиковать анонимно</label>
        </div>
        <button class="btn btn-primary" id="forum-submit">Опубликовать</button>
        <button class="btn btn-text" id="forum-cancel">Отмена</button>
      </div>
    `;
    const topicsEl = document.getElementById('forum-topics');
    const rootPosts = (list || []).filter(p => !p.parent_id);
    topicsEl.innerHTML = rootPosts.length
      ? rootPosts.map(p => `
        <div class="list-item forum-topic-click" data-post-id="${p.id}" style="cursor:pointer">
          <div>
            <strong>${escapeHtml(p.title || '(без заголовка)')}</strong>
            <p class="meta">${escapeHtml((p.body || '').slice(0, 120))}...</p>
            <p class="meta">${p.is_anonymous ? 'Анонимно' : (p.author_id ? `<span class="forum-author-click" data-user-id="${p.author_id}">${escapeHtml(p.author_name || '')}</span>` : escapeHtml(p.author_name || ''))} · ${formatDate(p.created_at)}</p>
          </div>
        </div>
      `).join('')
      : '<p class="empty-state">Нет тем</p>';

    topicsEl.querySelectorAll('.forum-topic-click').forEach((el, i) => {
      el.addEventListener('click', () => {
        const p = rootPosts[i];
        state.forumTopicId = p.id;
        state.forumTopic = p;
        loadForum();
      });
    });
    topicsEl.querySelectorAll('.forum-author-click').forEach(el => {
      el.style.cursor = 'pointer';
      el.addEventListener('click', (e) => { e.stopPropagation(); showUserProfile(el.dataset.userId); });
    });
    document.getElementById('forum-new-topic').addEventListener('click', () => {
      document.getElementById('forum-form').classList.remove('hidden');
    });
    document.getElementById('forum-cancel').addEventListener('click', () => {
      document.getElementById('forum-form').classList.add('hidden');
    });
    document.getElementById('forum-submit').addEventListener('click', () => {
      const title = document.getElementById('forum-title').value.trim();
      const body = document.getElementById('forum-body').value.trim();
      const isAnonymous = document.getElementById('forum-anonymous').checked;
      if (!body) { alert('Введите текст'); return; }
      api('/api/forum/posts', {
        method: 'POST',
        body: JSON.stringify({ parent_id: null, title: title || null, body, is_anonymous: isAnonymous })
      }).then(() => {
        document.getElementById('forum-form').classList.add('hidden');
        loadForum();
      }).catch(err => alert(err.message));
    });
  }).catch(err => {
    document.getElementById('content').innerHTML = `<p class="error-msg">${escapeHtml(err.message)}</p>`;
  });
}

function loadEvents() {
  const isAdmin = ['admin', 'director'].includes(state.user?.role || '');
  function fetchEvents(params) {
    const q = (params && params.q) ? '&q=' + encodeURIComponent(params.q) : '';
    const from = (params && params.from) ? '&from=' + encodeURIComponent(params.from) : '';
    const to = (params && params.to) ? '&to=' + encodeURIComponent(params.to) : '';
    return api('/api/events?limit=20' + q + from + to);
  }
  function renderEventsList(list) {
    const listEl = document.getElementById('events-list');
    if (!listEl) return;
    const events = Array.isArray(list) ? list : [];
    listEl.innerHTML = events.length
      ? events.map(e => `
        <div class="list-item">
          <div>
            <strong>${escapeHtml(e.title)}</strong>
            <p class="meta">${escapeHtml(e.description || '')}</p>
            <p class="meta">${formatDate(e.event_date)} ${e.location ? ' · ' + escapeHtml(e.location) : ''}</p>
          </div>
          ${isAdmin ? `<button class="btn btn-small btn-danger" data-ev-id="${e.id}">Удалить</button>` : ''}
        </div>
      `).join('')
      : '<p class="empty-state">Нет событий</p>';
    if (isAdmin) listEl.querySelectorAll('[data-ev-id]').forEach(btn => {
      btn.addEventListener('click', () => {
        if (!confirm('Удалить событие?')) return;
        api('/api/events/' + btn.dataset.evId, { method: 'DELETE' }).then(() => loadEvents()).catch(err => alert(err.message));
      });
    });
  }
  const formHtml = isAdmin ? `
    <div id="events-form" class="card hidden">
      <h3>Новое событие</h3>
      <div class="form-row"><label>Название</label><input type="text" id="ev-title"></div>
      <div class="form-row"><label>Описание</label><textarea id="ev-desc"></textarea></div>
      <div class="form-row"><label>Дата (ДД/ММ/ГГГГ)</label><input type="text" id="ev-date" placeholder="15/02/2026"></div>
      <div class="form-row"><label>Место</label><input type="text" id="ev-location"></div>
      <button class="btn btn-primary" id="ev-submit">Создать</button>
      <button class="btn btn-text" id="ev-cancel">Отмена</button>
    </div>
  ` : '';
  document.getElementById('content').innerHTML = `
    <div class="card">
      <h3>События и новости</h3>
      <div class="form-row" style="display:flex;flex-wrap:wrap;gap:0.5rem;margin-bottom:0.5rem">
        <input type="text" id="events-q" placeholder="Поиск по названию..." style="max-width:200px">
        <input type="text" id="events-from" placeholder="От (ДД/ММ/ГГГГ)">
        <input type="text" id="events-to" placeholder="До (ДД/ММ/ГГГГ)">
        <button type="button" class="btn btn-small btn-secondary" id="events-filter">Фильтр</button>
      </div>
      ${isAdmin ? '<button class="btn btn-primary btn-small" id="events-add">Добавить</button>' : ''}
      <div id="events-list" style="margin-top: 1rem"></div>
    </div>
    ${formHtml}
  `;
  fetchEvents().then(list => renderEventsList(list)).catch(err => {
    document.getElementById('content').innerHTML = `<p class="error-msg">${escapeHtml(err.message)}</p>`;
  });
  document.getElementById('events-filter').addEventListener('click', () => {
    const q = document.getElementById('events-q').value.trim();
    const fromRaw = document.getElementById('events-from').value.trim();
    const toRaw = document.getElementById('events-to').value.trim();
    const from = fromRaw ? parseDDMMYYYY(fromRaw) || fromRaw : undefined;
    const to = toRaw ? parseDDMMYYYY(toRaw) || toRaw : undefined;
    fetchEvents({ q: q || undefined, from, to }).then(renderEventsList).catch(err => alert(err.message));
  });
  if (isAdmin) {
    document.getElementById('events-add').addEventListener('click', () => document.getElementById('events-form').classList.remove('hidden'));
    document.getElementById('ev-cancel').addEventListener('click', () => document.getElementById('events-form').classList.add('hidden'));
    document.getElementById('ev-submit').addEventListener('click', () => {
      const title = document.getElementById('ev-title').value.trim();
      const dateRaw = document.getElementById('ev-date').value.trim();
      const date = parseDDMMYYYY(dateRaw) || dateRaw;
      if (!title || !date) { alert('Название и дата обязательны. Дата в формате ДД/ММ/ГГГГ'); return; }
      api('/api/events', {
        method: 'POST',
        body: JSON.stringify({
          title,
          description: document.getElementById('ev-desc').value.trim(),
          event_date: date + 'T12:00:00Z',
          location: document.getElementById('ev-location').value.trim()
        })
      }).then(() => { document.getElementById('events-form').classList.add('hidden'); loadEvents(); }).catch(err => alert(err.message));
    });
  }
}

function loadNotes() {
  if (state.notesNoteId) {
    api('/api/notes/' + state.notesNoteId).then(data => {
      const note = data.note || {};
      const comments = data.comments || [];
      document.getElementById('content').innerHTML = `
        <div class="card">
          <button type="button" class="btn btn-text btn-small" id="notes-back">← Назад к заметкам</button>
          <h3>${escapeHtml(note.title || 'Заметка')}</h3>
          <div class="forum-topic-body">${escapeHtml(note.body || '')}</div>
          <p class="meta">${formatDate(note.created_at)}</p>
          <h4 style="margin-top:1rem">Дополнения (${comments.length})</h4>
          <div id="notes-comments">${comments.map(c => `
            <div class="list-item"><div>${escapeHtml(c.body || '')}</div><p class="meta">${formatDate(c.created_at)}</p></div>
          `).join('')}</div>
          <div class="card" style="margin-top:1rem">
            <h4>Добавить дополнение</h4>
            <div class="form-row"><label>Текст</label><textarea id="notes-comment-body" placeholder="Текст дополнения"></textarea></div>
            <button type="button" class="btn btn-primary" id="notes-comment-submit">Добавить</button>
          </div>
        </div>
      `;
      document.getElementById('notes-back').addEventListener('click', () => {
        state.notesNoteId = null; state.notesNote = null; loadNotes();
      });
      document.getElementById('notes-comment-submit').addEventListener('click', () => {
        const body = document.getElementById('notes-comment-body').value.trim();
        if (!body) { alert('Введите текст'); return; }
        api('/api/notes/' + state.notesNoteId + '/comments', { method: 'POST', body: JSON.stringify({ body }) }).then(() => {
          document.getElementById('notes-comment-body').value = ''; loadNotes();
        }).catch(err => alert(err.message));
      });
    }).catch(err => { document.getElementById('content').innerHTML = `<p class="error-msg">${escapeHtml(err.message)}</p>`; });
    return;
  }
  state.notesNote = null;
  api('/api/notes').then(list => {
    const notes = Array.isArray(list) ? list : [];
    document.getElementById('content').innerHTML = `
      <div class="card">
        <h3>Мои заметки</h3>
        <button class="btn btn-primary btn-small" id="notes-new">Новая заметка</button>
        <div id="notes-topics" style="margin-top: 1rem"></div>
      </div>
      <div id="notes-form" class="card hidden">
        <h3>Новая заметка</h3>
        <div class="form-row">
          <label>Заголовок</label>
          <input type="text" id="notes-title" placeholder="Заголовок">
        </div>
        <div class="form-row">
          <label>Текст</label>
          <textarea id="notes-body" placeholder="Текст"></textarea>
        </div>
        <button class="btn btn-primary" id="notes-submit">Создать</button>
        <button class="btn btn-text" id="notes-cancel">Отмена</button>
      </div>
    `;
    const topicsEl = document.getElementById('notes-topics');
    topicsEl.innerHTML = notes.length
      ? notes.map(n => `
        <div class="list-item forum-topic-click" data-note-id="${n.id}" style="cursor:pointer">
          <div>
            <strong>${escapeHtml(n.title || '(без заголовка)')}</strong>
            <p class="meta">${escapeHtml((n.body || '').slice(0, 120))}${(n.body || '').length > 120 ? '...' : ''}</p>
            <p class="meta">${formatDate(n.created_at)}</p>
          </div>
        </div>
      `).join('')
      : '<p class="empty-state">Нет заметок</p>';

    topicsEl.querySelectorAll('.forum-topic-click').forEach((el, i) => {
      el.addEventListener('click', () => {
        const n = notes[i];
        state.notesNoteId = n.id;
        state.notesNote = n;
        loadNotes();
      });
    });
    document.getElementById('notes-new').addEventListener('click', () => {
      document.getElementById('notes-form').classList.remove('hidden');
    });
    document.getElementById('notes-cancel').addEventListener('click', () => {
      document.getElementById('notes-form').classList.add('hidden');
    });
    document.getElementById('notes-submit').addEventListener('click', () => {
      const title = document.getElementById('notes-title').value.trim();
      const body = document.getElementById('notes-body').value.trim();
      if (!body) { alert('Введите текст'); return; }
      api('/api/notes', {
        method: 'POST',
        body: JSON.stringify({ title: title || null, body })
      }).then(() => {
        document.getElementById('notes-form').classList.add('hidden');
        document.getElementById('notes-title').value = '';
        document.getElementById('notes-body').value = '';
        loadNotes();
      }).catch(err => alert(err.message));
    });
  }).catch(err => { document.getElementById('content').innerHTML = `<p class="error-msg">${escapeHtml(err.message)}</p>`; });
}

function loadChat() {
  document.getElementById('content').innerHTML = `
    <div class="card">
      <h3>Общий чат</h3>
      <div class="chat-messages" id="chat-messages"></div>
      <div class="chat-input-row">
        <input type="text" id="chat-input" placeholder="Сообщение...">
        <button class="btn btn-primary" id="chat-send">Отправить</button>
      </div>
    </div>
  `;
  const messagesEl = document.getElementById('chat-messages');
  const isAdmin = ['admin', 'director'].includes(state.user?.role || '');
  const myId = state.user?.user_id || state.user?.userId;

  function renderChat() {
    if (!messagesEl) return;
    messagesEl.innerHTML = (state.chatMessages || []).map(m => {
      let body = m.body || m.message || '';
      let author = m.author_name || '';
      let id = m.id;
      let userId = m.user_id;
      try { if (typeof m === 'string') m = JSON.parse(m); body = m.body || m.message || m; author = m.author_name || author; id = m.id; userId = m.user_id; } catch (_) {}
      if (typeof body !== 'string') body = String(body);
      const canDelete = id && (isAdmin || userId === myId);
      return `<div class="chat-msg" data-msg-id="${id || ''}">
        <div>${author ? `<span class="chat-author" data-user-id="${userId || ''}" style="cursor:pointer;font-weight:600">${escapeHtml(author)}:</span> ` : ''}${escapeHtml(body)}</div>
        ${canDelete ? `<button class="btn btn-text btn-small chat-msg-del" data-id="${id}">Удалить</button>` : ''}
      </div>`;
    }).join('');
    messagesEl.scrollTop = messagesEl.scrollHeight;
    messagesEl.querySelectorAll('.chat-author[data-user-id]').forEach(el => {
      const uid = el.dataset.userId;
      if (uid) el.addEventListener('click', () => showUserProfile(uid));
    });
    messagesEl.querySelectorAll('.chat-msg-del').forEach(btn => {
      btn.addEventListener('click', () => {
        api('/api/chat/messages/' + btn.dataset.id, { method: 'DELETE' }).then(() => {
          state.chatMessages = state.chatMessages.filter(m => m.id != btn.dataset.id);
          renderChat();
        }).catch(err => alert(err.message));
      });
    });
  }

  api('/api/chat/messages').then(list => {
    state.chatMessages = Array.isArray(list) ? list : [];
  }).catch(() => { state.chatMessages = []; }).finally(() => {
    if (!state.chatWs || state.chatWs.readyState === WebSocket.CLOSED || state.chatWs.readyState === WebSocket.CLOSING) {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const wsUrl = protocol + '//' + window.location.host + API_BASE + '/api/chat/ws?token=' + encodeURIComponent(state.token || '');
      const ws = new WebSocket(wsUrl);
      ws.onopen = () => {
        ws.send(JSON.stringify({ type: 'hello', body: state.user?.full_name + ' подключился' }));
      };
      ws.onmessage = (e) => {
        try {
          const m = JSON.parse(e.data);
          if (m.id && !state.chatMessages.find(x => x.id === m.id)) state.chatMessages.push(m);
          else if (!m.id) state.chatMessages.push(m);
        } catch (_) {
          state.chatMessages.push({ body: e.data });
        }
        renderChat();
      };
      ws.onerror = () => { state.chatMessages.push({ body: 'Ошибка подключения' }); renderChat(); };
      state.chatWs = ws;
    }
    renderChat();
  });

  document.getElementById('chat-send').addEventListener('click', () => {
    const input = document.getElementById('chat-input');
    const text = input.value.trim();
    if (!text || !state.chatWs || state.chatWs.readyState !== WebSocket.OPEN) return;
    const msg = { type: 'message', body: text };
    state.chatWs.send(JSON.stringify(msg));
    input.value = '';
  });
  document.getElementById('chat-input').addEventListener('keydown', (e) => {
    if (e.key === 'Enter') document.getElementById('chat-send').click();
  });
  renderChat();
}

function loadProfile() {
  api('/api/profile').then(p => {
    if (p) {
      state.user = { ...state.user, full_name: p.full_name, email: p.email, role: p.role, phone: p.phone, avatar_url: p.avatar_url };
      saveStored();
    }
    const u = state.user;
    document.getElementById('content').innerHTML = `
      <div class="card">
        <h3>Профиль</h3>
        <p><strong>Имя:</strong> ${escapeHtml(u?.full_name || '')}</p>
        <p><strong>Email:</strong> ${escapeHtml(u?.email || '')}</p>
        <p><strong>Телефон:</strong> ${escapeHtml(u?.phone || '—')}</p>
        <p><strong>Роль:</strong> ${escapeHtml(u?.role || '')}</p>
      </div>
      <div class="card">
        <h3>Сменить пароль и телефон</h3>
        <div class="form-row"><label>Новый пароль</label><input type="password" id="profile-password" placeholder="Новый пароль"></div>
        <div class="form-row"><label>Телефон</label><input type="text" id="profile-phone" placeholder="+7 ..." value="${escapeHtml(u?.phone || '')}"></div>
        <button type="button" class="btn btn-primary" id="profile-save">Сохранить</button>
      </div>
    `;
    document.getElementById('profile-save').addEventListener('click', () => {
      const password = document.getElementById('profile-password').value;
      const phone = document.getElementById('profile-phone').value.trim();
      const body = {};
      if (password) body.password = password;
      body.phone = phone;
      if (!password && !phone) { alert('Введите новый пароль и/или телефон'); return; }
      api('/api/profile', { method: 'PUT', body: JSON.stringify(body) }).then(() => {
        alert('Сохранено');
        document.getElementById('profile-password').value = '';
        state.user.phone = phone;
        saveStored();
        loadProfile();
      }).catch(err => alert(err.message));
    });
  }).catch(err => {
    const u = state.user;
    document.getElementById('content').innerHTML = `
      <div class="card">
        <h3>Профиль</h3>
        <p><strong>Имя:</strong> ${escapeHtml(u?.full_name || '')}</p>
        <p><strong>Email:</strong> ${escapeHtml(u?.email || '')}</p>
        <p><strong>Телефон:</strong> ${escapeHtml(u?.phone || '—')}</p>
        <p><strong>Роль:</strong> ${escapeHtml(u?.role || '')}</p>
      </div>
      <div class="card">
        <h3>Сменить пароль и телефон</h3>
        <div class="form-row"><label>Новый пароль</label><input type="password" id="profile-password" placeholder="Новый пароль"></div>
        <div class="form-row"><label>Телефон</label><input type="text" id="profile-phone" placeholder="+7 ..." value="${escapeHtml(u?.phone || '')}"></div>
        <button type="button" class="btn btn-primary" id="profile-save">Сохранить</button>
      </div>
    `;
    document.getElementById('profile-save').addEventListener('click', () => {
      const password = document.getElementById('profile-password').value;
      const phone = document.getElementById('profile-phone').value.trim();
      const body = {};
      if (password) body.password = password;
      body.phone = phone;
      if (!password && !phone) { alert('Введите новый пароль и/или телефон'); return; }
      api('/api/profile', { method: 'PUT', body: JSON.stringify(body) }).then(() => {
        alert('Сохранено');
        document.getElementById('profile-password').value = '';
        state.user.phone = phone;
        saveStored();
        loadProfile();
      }).catch(err => alert(err.message));
    });
  });
}

function loadDirectorProfile() {
  if (state.directorPostId) {
    api('/api/director/posts/' + state.directorPostId).then(data => {
      const post = data.post || {};
      const comments = data.comments || [];
      document.getElementById('content').innerHTML = `
        <div class="card">
          <button type="button" class="btn btn-text btn-small" id="director-back">← Назад к постам</button>
          <h3>${escapeHtml(post.title || 'Пост')}</h3>
          <div class="forum-topic-body">${escapeHtml(post.body || '')}</div>
          <p class="meta">${formatDate(post.created_at)}</p>
          <h4 style="margin-top:1rem">Комментарии (${comments.length})</h4>
          <div id="director-comments">${comments.map(c => `
            <div class="list-item"><div>${escapeHtml(c.body || '')}</div><p class="meta">${escapeHtml(c.author_name || '')} · ${formatDate(c.created_at)}</p></div>
          `).join('')}</div>
          <div class="card" style="margin-top:1rem">
            <h4>Комментировать</h4>
            <div class="form-row"><label>Текст</label><textarea id="director-comment-body" placeholder="Текст комментария"></textarea></div>
            <button type="button" class="btn btn-primary" id="director-comment-submit">Отправить</button>
          </div>
        </div>
      `;
      document.getElementById('director-back').addEventListener('click', () => {
        state.directorPostId = null; state.directorPost = null; loadDirectorProfile();
      });
      document.getElementById('director-comment-submit').addEventListener('click', () => {
        const body = document.getElementById('director-comment-body').value.trim();
        if (!body) { alert('Введите текст'); return; }
        api('/api/director/posts/' + state.directorPostId + '/comments', { method: 'POST', body: JSON.stringify({ body }) }).then(() => {
          document.getElementById('director-comment-body').value = ''; loadDirectorProfile();
        }).catch(err => alert(err.message));
      });
    }).catch(err => { document.getElementById('content').innerHTML = `<p class="error-msg">${escapeHtml(err.message)}</p>`; });
    return;
  }
  state.directorPost = null;
  Promise.all([
    fetch(API_BASE + '/api/director').then(async r => { if (!r.ok) return {}; try { return await r.json(); } catch (_) { return {}; } }),
    fetch(API_BASE + '/api/director/posts').then(async r => { if (!r.ok) return []; try { return await r.json(); } catch (_) { return []; } })
  ]).then(([d, postsList]) => {
    const posts = Array.isArray(postsList) ? postsList : [];
    const isDirector = state.user?.role === 'director';
    document.getElementById('content').innerHTML = `
      <div class="card director-card">
        <h3>Профиль директора</h3>
        <div class="director-header">
          <div><strong class="director-name">${escapeHtml(d.full_name || '')}</strong>
          ${d.phone ? `<p class="meta">Телефон: ${escapeHtml(d.phone)}</p>` : ''}
          ${d.email ? `<p class="meta">${escapeHtml(d.email)}</p>` : ''}
          <a href="tel:${escapeHtml(d.phone || '')}" class="btn btn-primary btn-small" style="margin-top:0.5rem">Написать директору</a></div>
        </div>
      </div>
      <div class="card">
        <h3>Посты директора</h3>
        ${isDirector ? '<button class="btn btn-primary btn-small" id="director-new-post">Новый пост</button>' : ''}
        <div id="director-posts-topics" style="margin-top: 1rem"></div>
      </div>
      ${isDirector ? `
      <div id="director-form" class="card hidden">
        <h3>Новый пост</h3>
        <div class="form-row">
          <label>Заголовок</label>
          <input type="text" id="director-post-title" placeholder="Заголовок">
        </div>
        <div class="form-row">
          <label>Текст</label>
          <textarea id="director-post-body" placeholder="Текст"></textarea>
        </div>
        <button class="btn btn-primary" id="director-post-submit">Опубликовать</button>
        <button class="btn btn-text" id="director-post-cancel">Отмена</button>
      </div>
      ` : ''}
    `;
    const topicsEl = document.getElementById('director-posts-topics');
    topicsEl.innerHTML = posts.length
      ? posts.map(p => `
        <div class="list-item forum-topic-click" data-post-id="${p.id}" style="cursor:pointer">
          <div>
            <strong>${escapeHtml(p.title || '(без заголовка)')}</strong>
            <p class="meta">${escapeHtml((p.body || '').slice(0, 120))}${(p.body || '').length > 120 ? '...' : ''}</p>
            <p class="meta">${formatDate(p.created_at)}</p>
          </div>
        </div>
      `).join('')
      : '<p class="empty-state">Нет постов</p>';

    topicsEl.querySelectorAll('.forum-topic-click').forEach((el, i) => {
      el.addEventListener('click', () => {
        const p = posts[i];
        state.directorPostId = p.id;
        state.directorPost = p;
        loadDirectorProfile();
      });
    });
    if (isDirector) {
      document.getElementById('director-new-post').addEventListener('click', () => {
        document.getElementById('director-form').classList.remove('hidden');
      });
      document.getElementById('director-post-cancel').addEventListener('click', () => {
        document.getElementById('director-form').classList.add('hidden');
      });
      document.getElementById('director-post-submit').addEventListener('click', () => {
        const title = document.getElementById('director-post-title').value.trim();
        const body = document.getElementById('director-post-body').value.trim();
        if (!body) { alert('Введите текст'); return; }
        api('/api/director/posts', {
          method: 'POST',
          body: JSON.stringify({ title: title || null, body })
        }).then(() => {
          document.getElementById('director-form').classList.add('hidden');
          document.getElementById('director-post-title').value = '';
          document.getElementById('director-post-body').value = '';
          loadDirectorProfile();
        }).catch(err => alert(err.message));
      });
    }
  }).catch(err => { document.getElementById('content').innerHTML = `<p class="error-msg">${escapeHtml(err.message)}</p>`; });
}

function loadRegulations() {
  document.getElementById('content').innerHTML = `
    <div class="card">
      <h3>Регламент и академический календарь</h3>
      <p>Академический календарь, правила обучения и регламент размещены на официальном сайте колледжа.</p>
      <p class="meta">Силлабусы и материалы по дисциплинам доступны в разделе «Расписание» — кнопка «ДЗ и материалы» у каждого занятия.</p>
      <p><a href="https://collegenarxoz.kz/" target="_blank" rel="noopener" class="btn btn-secondary btn-small">Сайт колледжа НАРХОЗ</a></p>
    </div>
  `;
}

function loadNotifications() {
  const isAdmin = ['admin', 'director'].includes(state.user?.role || '');
  api('/api/notifications?limit=50').then(list => {
    const items = Array.isArray(list) ? list : [];
    document.getElementById('content').innerHTML = `
      <div class="card">
        <h3>Уведомления</h3>
        ${isAdmin ? `
          <div class="form-row">
            <input type="text" id="notif-title" placeholder="Заголовок">
            <textarea id="notif-body" placeholder="Текст уведомления" rows="2"></textarea>
            <button class="btn btn-primary btn-small" id="notif-send">Отправить всем</button>
          </div>
        ` : ''}
        <div id="notif-list" style="margin-top:1rem">
          ${items.length ? items.map(n => `
            <div class="list-item">
              <div>
                <strong>${escapeHtml(n.title)}</strong>
                ${n.body ? `<p class="meta">${escapeHtml(n.body)}</p>` : ''}
                <p class="meta">${formatDate(n.created_at)}</p>
              </div>
            </div>
          `).join('') : '<p class="empty-state">Нет уведомлений</p>'}
        </div>
      </div>
    `;
    if (isAdmin) {
      document.getElementById('notif-send').addEventListener('click', () => {
        const title = document.getElementById('notif-title').value.trim();
        const body = document.getElementById('notif-body').value.trim();
        if (!title) { alert('Введите заголовок'); return; }
        api('/api/notifications', { method: 'POST', body: JSON.stringify({ title, body }) }).then(() => {
          document.getElementById('notif-title').value = '';
          document.getElementById('notif-body').value = '';
          loadNotifications();
        }).catch(err => alert(err.message));
      });
    }
  }).catch(err => {
    document.getElementById('content').innerHTML = `<p class="error-msg">${escapeHtml(err.message)}</p>`;
  });
}

function loadAllUsers() {
  api('/api/users').then(list => {
    const items = Array.isArray(list) ? list : [];
    document.getElementById('content').innerHTML = `
      <div class="card">
        <h3>Все пользователи</h3>
        <p class="meta">Имя · Роль</p>
        <div id="users-list">
          ${items.length ? items.map(u => `
            <div class="list-item user-item" data-user-id="${u.id}">
              <div>
                <strong>${escapeHtml(u.full_name)}</strong>
                <span class="meta"> · ${escapeHtml(u.role)}</span>
                ${u.email ? `<p class="meta">${escapeHtml(u.email)}</p>` : ''}
              </div>
              <div style="display:flex;gap:0.5rem">
                <button class="btn btn-small btn-secondary btn-edit-user" data-user-id="${u.id}">Редактировать</button>
                <button class="btn btn-small btn-danger btn-delete-user" data-user-id="${u.id}">Удалить</button>
              </div>
            </div>
          `).join('') : '<p class="empty-state">Нет пользователей</p>'}
        </div>
      </div>
      <div id="user-edit-modal" class="modal hidden">
        <div class="modal-content card">
          <h3>Редактировать пользователя</h3>
          <input type="hidden" id="edit-user-id">
          <div class="form-row"><label>Имя</label><input type="text" id="edit-user-name"></div>
          <div class="form-row"><label>Email</label><input type="email" id="edit-user-email"></div>
          <div class="form-row"><label>Телефон</label><input type="text" id="edit-user-phone"></div>
          <div class="form-row"><label>Роль</label><select id="edit-user-role"><option value="student">student</option><option value="teacher">teacher</option><option value="director">director</option><option value="admin">admin</option></select></div>
          <button class="btn btn-primary" id="edit-user-save">Сохранить</button>
          <button class="btn btn-text" id="edit-user-cancel">Отмена</button>
        </div>
      </div>
    `;
    const users = items;
    document.querySelectorAll('.btn-edit-user').forEach(btn => {
      btn.addEventListener('click', () => {
        const id = btn.dataset.userId;
        const u = users.find(x => String(x.id) === id);
        if (!u) return;
        document.getElementById('edit-user-id').value = u.id;
        document.getElementById('edit-user-name').value = u.full_name || '';
        document.getElementById('edit-user-email').value = u.email || '';
        document.getElementById('edit-user-phone').value = u.phone || '';
        document.getElementById('edit-user-role').value = u.role || 'student';
        document.getElementById('user-edit-modal').classList.remove('hidden');
      });
    });
    document.querySelectorAll('.btn-delete-user').forEach(btn => {
      btn.addEventListener('click', () => {
        if (!confirm('Деактивировать пользователя?')) return;
        api('/api/users/' + btn.dataset.userId, { method: 'DELETE' }).then(() => loadAllUsers()).catch(err => alert(err.message));
      });
    });
    document.getElementById('edit-user-save').addEventListener('click', () => {
      const id = document.getElementById('edit-user-id').value;
      api('/api/users/' + id, {
        method: 'PUT',
        body: JSON.stringify({
          full_name: document.getElementById('edit-user-name').value.trim(),
          email: document.getElementById('edit-user-email').value.trim(),
          phone: document.getElementById('edit-user-phone').value.trim(),
          role: document.getElementById('edit-user-role').value
        })
      }).then(() => {
        document.getElementById('user-edit-modal').classList.add('hidden');
        loadAllUsers();
      }).catch(err => alert(err.message));
    });
    document.getElementById('edit-user-cancel').addEventListener('click', () => {
      document.getElementById('user-edit-modal').classList.add('hidden');
    });
  }).catch(err => {
    document.getElementById('content').innerHTML = `<p class="error-msg">${escapeHtml(err.message)}</p>`;
  });
}

function loadMap() {
  const isAdmin = ['admin', 'director'].includes(state.user?.role || '');
  fetch(API_BASE + '/api/building-map').then(r => r.ok ? r.json() : {}).then(data => {
    const content = data?.content || 'Карта и план эвакуации. Размещаются администрацией.';
    const imageUrl = data?.image_url || '';
    document.getElementById('content').innerHTML = `
      <div class="card">
        <h3>Карта здания</h3>
        <p class="meta">Адрес: 10-й микрорайон, 7а/2, г. Алматы, Казахстан, 050035.</p>
        <div class="map-placeholder">
          ${imageUrl ? `<img src="${API_BASE}${imageUrl}" alt="Карта" style="max-width:100%;border-radius:2px">` : ''}
          <p>${escapeHtml(content)}</p>
        </div>
        ${isAdmin ? `
          <div class="card" style="margin-top:1rem">
            <h4>Редактировать</h4>
            <textarea id="map-content" rows="3" placeholder="Описание">${escapeHtml(content)}</textarea>
            <input type="text" id="map-image-url" placeholder="URL изображения карты" value="${escapeHtml(imageUrl)}">
            <button class="btn btn-primary btn-small" id="map-save">Сохранить</button>
          </div>
        ` : ''}
      </div>
    `;
    if (isAdmin) {
      document.getElementById('map-save').addEventListener('click', () => {
        api('/api/building-map', {
          method: 'PUT',
          body: JSON.stringify({
            content: document.getElementById('map-content').value.trim(),
            image_url: document.getElementById('map-image-url').value.trim()
          })
        }).then(() => { alert('Сохранено'); loadMap(); }).catch(err => alert(err.message));
      });
    }
  }).catch(() => {
    document.getElementById('content').innerHTML = `
      <div class="card">
        <h3>Карта здания</h3>
        <p class="meta">Адрес: 10-й микрорайон, 7а/2, г. Алматы, Казахстан, 050035.</p>
        <div class="map-placeholder"><span>Карта и план эвакуации</span></div>
      </div>
    `;
  });
}

// ——— Вход / выход ———
function logout() {
  state.token = null;
  state.user = null;
  clearStored();
  if (state.chatWs) {
    state.chatWs.close();
    state.chatWs = null;
  }
  document.getElementById('screen-login').classList.add('active');
  document.getElementById('screen-app').classList.remove('active');
}

function showApp() {
  document.getElementById('screen-login').classList.remove('active');
  document.getElementById('screen-app').classList.add('active');
  document.getElementById('user-name').textContent = state.user?.full_name || '';
  document.getElementById('user-role').textContent = state.user?.role || '';
  renderNav();
  goTo(state.currentPage);
}

// ——— Инициализация ———
function showUserProfile(userId) {
  if (!userId) return;
  api('/api/users/' + userId).then(u => {
    const modal = document.createElement('div');
    modal.className = 'modal';
    modal.innerHTML = `
      <div class="modal-content card">
        <h3>Профиль</h3>
        ${u.avatar_url ? `<img src="${API_BASE}${u.avatar_url}" alt="" style="width:64px;height:64px;border-radius:2px;object-fit:cover;margin-bottom:0.5rem">` : ''}
        <p><strong>${escapeHtml(u.full_name || '')}</strong></p>
        <p class="meta">${escapeHtml(u.email || '')}</p>
        <p class="meta">${escapeHtml(u.phone || '')}</p>
        <p class="meta">${escapeHtml(u.role || '')}</p>
        <button class="btn btn-text" id="user-profile-close">Закрыть</button>
      </div>
    `;
    modal.querySelector('#user-profile-close').onclick = () => modal.remove();
    modal.onclick = (e) => { if (e.target === modal) modal.remove(); };
    document.body.appendChild(modal);
  }).catch(() => {});
}

function escapeHtml(s) {
  if (s == null) return '';
  const div = document.createElement('div');
  div.textContent = s;
  return div.innerHTML;
}

function formatDate(s) {
  if (!s) return '';
  const d = new Date(s);
  if (isNaN(d.getTime())) return s;
  const day = String(d.getDate()).padStart(2, '0');
  const month = String(d.getMonth() + 1).padStart(2, '0');
  const year = d.getFullYear();
  return day + '/' + month + '/' + year;
}

function toISOForInput(s) {
  if (!s) return '';
  const d = new Date(s);
  return isNaN(d.getTime()) ? '' : d.toISOString().slice(0, 10);
}

function parseDDMMYYYY(s) {
  if (!s || typeof s !== 'string') return '';
  const parts = s.trim().split(/[/.-]/);
  if (parts.length !== 3) return '';
  const day = parseInt(parts[0], 10);
  const month = parseInt(parts[1], 10) - 1;
  const year = parseInt(parts[2], 10);
  if (isNaN(day) || isNaN(month) || isNaN(year)) return '';
  const d = new Date(year, month, day);
  if (isNaN(d.getTime())) return '';
  return d.getFullYear() + '-' + String(d.getMonth() + 1).padStart(2, '0') + '-' + String(d.getDate()).padStart(2, '0');
}

function formatTime(d) {
  return d.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' });
}

document.getElementById('form-login').addEventListener('submit', async (e) => {
  e.preventDefault();
  const errEl = document.getElementById('login-error');
  errEl.classList.add('hidden');
  const email = e.target.email.value.trim();
  const password = e.target.password.value;
  if (!email || !password) {
    errEl.textContent = 'Введите email и пароль';
    errEl.classList.remove('hidden');
    return;
  }
  try {
    const res = await api('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password })
    });
    state.token = res.token;
    state.user = { user_id: res.user_id, role: res.role, email: res.email, full_name: res.full_name };
    saveStored();
    showApp();
  } catch (err) {
    errEl.textContent = err.message || 'Ошибка входа';
    errEl.classList.remove('hidden');
  }
});

document.getElementById('btn-logout').addEventListener('click', () => {
  logout();
});

document.getElementById('btn-menu').addEventListener('click', () => {
  var sidebar = document.getElementById('sidebar');
  var overlay = document.getElementById('sidebar-overlay');
  sidebar.classList.toggle('open');
  if (overlay) {
    overlay.classList.toggle('visible', sidebar.classList.contains('open'));
    overlay.classList.toggle('hidden', !sidebar.classList.contains('open'));
  }
});
var overlay = document.getElementById('sidebar-overlay');
if (overlay) overlay.addEventListener('click', function () {
  document.getElementById('sidebar').classList.remove('open');
  overlay.classList.add('hidden');
  overlay.classList.remove('visible');
});

if (loadStored()) {
  showApp();
} else {
  document.getElementById('screen-login').classList.add('active');
  document.getElementById('screen-app').classList.remove('active');
}
