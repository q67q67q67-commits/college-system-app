import 'package:flutter/material.dart';
import '../services/api_service.dart';

class TeacherHomeworkScreen extends StatefulWidget {
  const TeacherHomeworkScreen({super.key});

  @override
  State<TeacherHomeworkScreen> createState() => _TeacherHomeworkScreenState();
}

class _TeacherHomeworkScreenState extends State<TeacherHomeworkScreen> {
  List<dynamic> _schedule = [];
  Map<int, List<dynamic>> _attachments = {};
  String? _error;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final res = await ApiService.get('/api/schedule');
      setState(() {
        _schedule = res is List ? res : [];
        _error = null;
      });
      for (final s in _schedule) {
        _loadAttachments((s['id'] as num).toInt());
      }
    } catch (e) {
      setState(() {
        _schedule = [];
        _error = e.toString();
      });
    }
  }

  Future<void> _loadAttachments(int scheduleId) async {
    try {
      final res = await ApiService.get('/api/schedule/$scheduleId/attachments');
      setState(() => _attachments[scheduleId] = res is List ? res : []);
    } catch (_) {
      setState(() => _attachments[scheduleId] = []);
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_error != null) return Center(child: Column(mainAxisSize: MainAxisSize.min, children: [Text(_error!, style: const TextStyle(color: Colors.red)), TextButton(onPressed: _load, child: const Text('Повторить'))]));
    return RefreshIndicator(
      onRefresh: _load,
      child: ListView.builder(
        padding: const EdgeInsets.all(16),
        itemCount: _schedule.length,
        itemBuilder: (_, i) {
          final s = _schedule[i] as Map<String, dynamic>;
          final id = (s['id'] as num).toInt();
          final att = _attachments[id] ?? [];
          return Card(
            margin: const EdgeInsets.only(bottom: 8),
            child: ExpansionTile(
              title: Text(s['subject']?.toString() ?? ''),
              subtitle: Text('${s['room'] ?? ''} · ДЗ/материалов: ${att.length}'),
              children: [
                ...att.map<Widget>((a) => ListTile(
                      title: Text(a['title']?.toString() ?? ''),
                      subtitle: Text(() { final b = a['body']?.toString() ?? ''; return b.length > 60 ? '${b.substring(0, 60)}...' : b; }()),
                      trailing: IconButton(
                        icon: const Icon(Icons.delete),
                        onPressed: () async {
                          try {
                            await ApiService.delete('/api/attachments/${a['id']}');
                            if (mounted) _loadAttachments(id);
                          } catch (_) {}
                        },
                      ),
                    )),
                ListTile(
                  leading: const Icon(Icons.add),
                  title: const Text('Добавить ДЗ'),
                  onTap: () => _showAddDialog(id),
                ),
              ],
            ),
          );
        },
      ),
    );
  }

  Future<void> _showAddDialog(int scheduleId) async {
    final title = TextEditingController();
    final body = TextEditingController();
    await showDialog(
      context: context,
      builder: (_) => AlertDialog(
        title: const Text('Добавить ДЗ'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(controller: title, decoration: const InputDecoration(labelText: 'Название', border: OutlineInputBorder())),
            const SizedBox(height: 8),
            TextField(controller: body, decoration: const InputDecoration(labelText: 'Текст', border: OutlineInputBorder()), maxLines: 3),
          ],
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(context), child: const Text('Отмена')),
          FilledButton(
            onPressed: () async {
              if (title.text.trim().isEmpty) return;
              try {
                await ApiService.post('/api/schedule/$scheduleId/attachments', {'title': title.text.trim(), 'body': body.text.trim(), 'attachment_type': 'homework'});
                if (context.mounted) {
                  Navigator.pop(context);
                  _loadAttachments(scheduleId);
                }
              } catch (_) {}
            },
            child: const Text('Сохранить'),
          ),
        ],
      ),
    );
  }
}
