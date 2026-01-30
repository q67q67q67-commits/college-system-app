/**
 * Narxoz College (NC) — веб-прототип
 * Палитра: #a82523, #d50032 (сайт университета)
 */

const API_BASE = '';

let state = {
  token: null,
  user: null,
  currentPage: 'home',
  chatWs: null,
  chatMessages: [],
  forumTopicId: null,
  forumTopic: null
};

// ——— API ———
function getAuthHeaders() {
  const h = { 'Content-Type': 'application/json' };
  if (state.token) h['Authorization'] = 'Bearer ' + state.token;
  return h;
}

async function api(path, options = {}) {
  const url = path.startsWith('http') ? path : API_BASE + path;
  const res = await fetch(url, {
    ...options,
    headers: { ...getAuthHeaders(), ...(options.headers || {}) }
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
  { id: 'home', title: 'Главная', icon: '🏠', roles: ['student', 'teacher', 'director', 'admin'] },
  { id: 'schedule', title: 'Расписание', icon: '📅', roles: ['student', 'teacher', 'director', 'admin'] },
  { id: 'grades', title: 'Оценки', icon: '📊', roles: ['student'] },
  { id: 'groups', title: 'Группы', icon: '👥', roles: ['teacher', 'director', 'admin'] },
  { id: 'teacherGrades', title: 'Выставить оценку', icon: '✏️', roles: ['teacher', 'director', 'admin'] },
  { id: 'teacherHomework', title: 'ДЗ к парам', icon: '📋', roles: ['teacher', 'director', 'admin'] },
  { id: 'library', title: 'Библиотека', icon: '📚', roles: ['student', 'teacher', 'director', 'admin'] },
  { id: 'forum', title: 'Форум', icon: '💬', roles: ['student', 'teacher', 'director', 'admin'] },
  { id: 'events', title: 'События', icon: '📌', roles: ['student', 'teacher', 'director', 'admin'] },
  { id: 'files', title: 'Мои файлы', icon: '📁', roles: ['student', 'teacher', 'director', 'admin'] },
  { id: 'chat', title: 'Чат', icon: '💭', roles: ['student', 'teacher', 'director', 'admin'] },
  { id: 'directorProfile', title: 'Профиль директора', icon: '👔', roles: ['student', 'teacher', 'director'] },
  { id: 'regulations', title: 'Регламент', icon: '📜', roles: ['student', 'teacher', 'director', 'admin'] },
  { id: 'map', title: 'Карта здания', icon: '🗺️', roles: ['student', 'teacher', 'director', 'admin'] },
  { id: 'profile', title: 'Профиль', icon: '👤', roles: ['student', 'teacher', 'director', 'admin'] }
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
        <span>${p.icon}</span>
        <span style="margin-left: 0.5rem">${p.title}</span>
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
        <span>${p.icon}</span>
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
  else if (pageId === 'library') loadLibrary();
  else if (pageId === 'forum') loadForum();
  else if (pageId === 'events') loadEvents();
  else if (pageId === 'files') loadFiles();
  else if (pageId === 'chat') loadChat();
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
        <p><strong>GPA:</strong> ${gpa}</p>
        ${Array.isArray(list) && list.length
          ? list.map(g => `
            <div class="list-item">
              <div>
                <strong>${escapeHtml(g.subject)}</strong>
                <p class="meta">${g.grade_date} — ${g.grade} ${g.comment ? ' · ' + escapeHtml(g.comment) : ''}</p>
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
          <label>Дата (YYYY-MM-DD)</label>
          <input type="date" id="tg-date" value="${new Date().toISOString().slice(0, 10)}">
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
      const gradeDate = document.getElementById('tg-date').value;
      const comment = document.getElementById('tg-comment').value.trim();
      if (!userId || !scheduleId || isNaN(grade)) { alert('Заполните группу, студента, занятие и оценку'); return; }
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
              <button class="btn btn-small btn-secondary" data-group-id="${g.id}">Студенты</button>
            </div>
          `).join('')
          : '<p class="empty-state">Нет групп</p>'}
      </div>
    `;
    document.getElementById('content').querySelectorAll('[data-group-id]').forEach(btn => {
      btn.addEventListener('click', () => {
        const id = btn.dataset.groupId;
        api('/api/groups/' + id + '/students').then(students => {
          alert(students.map(s => s.full_name + ' (' + s.email + ')').join('\n') || 'Нет студентов');
        });
      });
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
          btn.addEventListener('click', () => reserveBook(btn.dataset.bookId));
        });
      });
    });
    document.getElementById('content').querySelectorAll('[data-book-id]').forEach(btn => {
      btn.addEventListener('click', () => reserveBook(btn.dataset.bookId));
    });
  }).catch(err => {
    document.getElementById('content').innerHTML = `<p class="error-msg">${escapeHtml(err.message)}</p>`;
  });
}

function reserveBook(bookId) {
  api('/api/library/reservations', {
    method: 'POST',
    body: JSON.stringify({ book_id: parseInt(bookId, 10) })
  }).then(() => {
    alert('Бронь создана');
    loadLibrary();
  }).catch(err => alert(err.message));
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
            <div class="list-item">
              <div>${escapeHtml(r.body || '')}</div>
              <p class="meta">${r.is_anonymous ? 'Анонимно' : (r.author_name || '')} · ${formatDate(r.created_at)}</p>
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
            <p class="meta">${p.is_anonymous ? 'Анонимно' : (p.author_name || '')} · ${formatDate(p.created_at)}</p>
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
      <div class="form-row"><label>Дата (YYYY-MM-DD)</label><input type="text" id="ev-date" placeholder="2026-02-15"></div>
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
        <input type="date" id="events-from" placeholder="От">
        <input type="date" id="events-to" placeholder="До">
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
    const from = document.getElementById('events-from').value;
    const to = document.getElementById('events-to').value;
    fetchEvents({ q: q || undefined, from: from || undefined, to: to || undefined }).then(renderEventsList).catch(err => alert(err.message));
  });
  if (isAdmin) {
    document.getElementById('events-add').addEventListener('click', () => document.getElementById('events-form').classList.remove('hidden'));
    document.getElementById('ev-cancel').addEventListener('click', () => document.getElementById('events-form').classList.add('hidden'));
    document.getElementById('ev-submit').addEventListener('click', () => {
      const title = document.getElementById('ev-title').value.trim();
      const date = document.getElementById('ev-date').value.trim();
      if (!title || !date) { alert('Название и дата обязательны'); return; }
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

function loadFiles() {
  api('/api/files').then(data => {
    const files = data.files || [];
    const used = data.storage_used ?? 0;
    const limit = data.storage_limit ?? 2 * 1024 * 1024 * 1024;
    const usedMB = (used / 1024 / 1024).toFixed(2);
    const limitGB = (limit / 1024 / 1024 / 1024).toFixed(1);
    document.getElementById('content').innerHTML = `
      <div class="card">
        <h3>Мои файлы</h3>
        <p class="meta">Использовано: ${usedMB} МБ из ${limitGB} ГБ</p>
        <div class="form-row">
          <input type="file" id="file-upload-input">
          <button class="btn btn-primary btn-small" id="file-upload-btn">Загрузить</button>
        </div>
        <div id="files-list" style="margin-top: 1rem"></div>
      </div>
    `;
    const listEl = document.getElementById('files-list');
    listEl.innerHTML = files.length
      ? files.map(f => `
        <div class="list-item">
          <div>
            <strong>${escapeHtml(f.filename || f.path)}</strong>
            <p class="meta">${(f.size_bytes / 1024).toFixed(1)} КБ</p>
          </div>
          <button class="btn btn-small btn-danger" data-file-id="${f.id}">Удалить</button>
        </div>
      `).join('')
      : '<p class="empty-state">Нет файлов</p>';

    document.getElementById('file-upload-btn').addEventListener('click', () => {
      const input = document.getElementById('file-upload-input');
      if (!input.files || !input.files[0]) { alert('Выберите файл'); return; }
      const fd = new FormData();
      fd.append('file', input.files[0]);
      fetch(API_BASE + '/api/files/upload', {
        method: 'POST',
        headers: { 'Authorization': 'Bearer ' + state.token },
        body: fd
      }).then(r => {
        if (r.status === 401) { logout(); window.location.reload(); return; }
        return r.json();
      }).then(data => {
        if (data && data.id) { loadFiles(); input.value = ''; }
        else alert(data?.error || 'Ошибка загрузки');
      }).catch(() => alert('Ошибка загрузки'));
    });
    listEl.querySelectorAll('[data-file-id]').forEach(btn => {
      btn.addEventListener('click', () => {
        if (!confirm('Удалить файл?')) return;
        api('/api/files/' + btn.dataset.fileId, { method: 'DELETE' }).then(() => loadFiles()).catch(err => alert(err.message));
      });
    });
  }).catch(err => {
    document.getElementById('content').innerHTML = `<p class="error-msg">${escapeHtml(err.message)}</p>`;
  });
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
  state.chatMessages = [];

  function renderChat() {
    messagesEl.innerHTML = state.chatMessages.map(m => {
      let body = m.body || m.message || '';
      try { if (typeof m === 'string') m = JSON.parse(m); body = m.body || m.message || m; } catch (_) {}
      if (typeof body !== 'string') body = String(body);
      return `<div class="chat-msg"><div>${escapeHtml(body)}</div></div>`;
    }).join('');
    messagesEl.scrollTop = messagesEl.scrollHeight;
  }

  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  const wsUrl = protocol + '//' + window.location.host + API_BASE + '/api/chat/ws?token=' + encodeURIComponent(state.token || '');
  const ws = new WebSocket(wsUrl);
  ws.onopen = () => {
    ws.send(JSON.stringify({ type: 'hello', body: state.user?.full_name + ' подключился' }));
  };
  ws.onmessage = (e) => {
    try {
      const m = JSON.parse(e.data);
      state.chatMessages.push(m);
    } catch (_) {
      state.chatMessages.push({ body: e.data });
    }
    renderChat();
  };
  ws.onerror = () => { state.chatMessages.push({ body: 'Ошибка подключения' }); renderChat(); };
  state.chatWs = ws;

  document.getElementById('chat-send').addEventListener('click', () => {
    const input = document.getElementById('chat-input');
    const text = input.value.trim();
    if (!text || !state.chatWs || state.chatWs.readyState !== WebSocket.OPEN) return;
    const msg = { type: 'message', body: text };
    state.chatWs.send(JSON.stringify(msg));
    state.chatMessages.push(msg);
    renderChat();
    input.value = '';
  });
  document.getElementById('chat-input').addEventListener('keydown', (e) => {
    if (e.key === 'Enter') document.getElementById('chat-send').click();
  });
  renderChat();
}

function loadProfile() {
  const u = state.user;
  document.getElementById('content').innerHTML = `
    <div class="card">
      <h3>Профиль</h3>
      <p><strong>Имя:</strong> ${escapeHtml(u?.full_name || '')}</p>
      <p><strong>Email:</strong> ${escapeHtml(u?.email || '')}</p>
      <p><strong>Роль:</strong> ${escapeHtml(u?.role || '')}</p>
    </div>
    <div class="card">
      <h3>Сменить пароль</h3>
      <div class="form-row"><label>Новый пароль</label><input type="password" id="profile-password" placeholder="Новый пароль"></div>
      <div class="form-row"><label>Телефон (опционально)</label><input type="text" id="profile-phone" placeholder="+7 ..."></div>
      <button type="button" class="btn btn-primary" id="profile-save">Сохранить</button>
    </div>
  `;
  document.getElementById('profile-save').addEventListener('click', () => {
    const password = document.getElementById('profile-password').value;
    const phone = document.getElementById('profile-phone').value.trim();
    const body = {};
    if (password) body.password = password;
    if (phone) body.phone = phone;
    if (!password && !phone) { alert('Введите новый пароль и/или телефон'); return; }
    api('/api/profile', { method: 'PUT', body: JSON.stringify(body) }).then(() => {
      alert('Сохранено');
      document.getElementById('profile-password').value = '';
    }).catch(err => alert(err.message));
  });
}

function loadDirectorProfile() {
  api('/api/director').then(d => {
    document.getElementById('content').innerHTML = `
      <div class="card director-card">
        <h3>Профиль директора</h3>
        <div class="director-header">
          ${d.avatar_url ? `<img src="${escapeHtml(d.avatar_url)}" alt="" class="director-avatar">` : '<div class="director-avatar-placeholder">👔</div>'}
          <div>
            <strong class="director-name">${escapeHtml(d.full_name || '')}</strong>
            ${d.phone ? `<p class="meta">Телефон: ${escapeHtml(d.phone)}</p>` : ''}
            ${d.email ? `<p class="meta">${escapeHtml(d.email)}</p>` : ''}
            <a href="tel:${escapeHtml(d.phone || '')}" class="btn btn-primary btn-small" style="margin-top:0.5rem">Написать директору</a>
          </div>
        </div>
      </div>
    `;
  }).catch(err => {
    document.getElementById('content').innerHTML = `<p class="error-msg">${escapeHtml(err.message)}</p>`;
  });
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

function loadMap() {
  document.getElementById('content').innerHTML = `
    <div class="card">
      <h3>Карта здания</h3>
      <p>Схема этажей и навигация по корпусу. При необходимости обратитесь в приёмную комиссию.</p>
      <p class="meta">Адрес: ул. Жандосова, 55, г. Алматы, Казахстан, 050035.</p>
      <div class="map-placeholder"><span>Карта и план эвакуации</span><br><small>Размещаются администрацией</small></div>
    </div>
  `;
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
function escapeHtml(s) {
  if (s == null) return '';
  const div = document.createElement('div');
  div.textContent = s;
  return div.innerHTML;
}

function formatDate(s) {
  if (!s) return '';
  const d = new Date(s);
  return isNaN(d.getTime()) ? s : d.toLocaleDateString('ru-RU');
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
