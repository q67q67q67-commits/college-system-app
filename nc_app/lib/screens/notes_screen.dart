import 'package:flutter/material.dart';
import '../services/api_service.dart';

class NotesScreen extends StatefulWidget {
  const NotesScreen({super.key});

  @override
  State<NotesScreen> createState() => _NotesScreenState();
}

class _NotesScreenState extends State<NotesScreen> {
  List<dynamic> _list = [];
  int? _selectedNoteId;
  Map<String, dynamic>? _noteDetail;
  List<dynamic> _comments = [];
  String? _error;
  bool _showForm = false;
  final _titleController = TextEditingController();
  final _bodyController = TextEditingController();
  final _commentController = TextEditingController();

  @override
  void dispose() {
    _titleController.dispose();
    _bodyController.dispose();
    _commentController.dispose();
    super.dispose();
  }

  Future<void> _loadList() async {
    try {
      final res = await ApiService.get('/api/notes');
      setState(() {
        _list = res is List ? List<dynamic>.from(res as List) : [];
        _error = null;
      });
    } catch (e) {
      setState(() {
        _list = [];
        _error = e.toString();
      });
    }
  }

  Future<void> _loadNote(int id) async {
    try {
      final res = await ApiService.get('/api/notes/$id');
      setState(() {
        _noteDetail = res['note'] is Map ? Map<String, dynamic>.from(res['note'] as Map) : null;
        _comments = res['comments'] is List ? List<dynamic>.from(res['comments'] as List) : [];
        _error = null;
      });
    } catch (e) {
      setState(() {
        _noteDetail = null;
        _comments = [];
        _error = e.toString();
      });
    }
  }

  Future<void> _createNote() async {
    final title = _titleController.text.trim();
    final body = _bodyController.text.trim();
    if (body.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Введите текст')));
      return;
    }
    try {
      await ApiService.post('/api/notes', {'title': title, 'body': body});
      _titleController.clear();
      _bodyController.clear();
      setState(() => _showForm = false);
      _loadList();
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(e.toString())));
    }
  }

  Future<void> _addComment() async {
    final body = _commentController.text.trim();
    if (body.isEmpty || _selectedNoteId == null) return;
    try {
      await ApiService.post('/api/notes/$_selectedNoteId/comments', {'body': body});
      _commentController.clear();
      _loadNote(_selectedNoteId!);
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(e.toString())));
    }
  }

  void _goBack() {
    setState(() {
      _selectedNoteId = null;
      _noteDetail = null;
      _comments = [];
    });
    _loadList();
  }

  @override
  void initState() {
    super.initState();
    _loadList();
  }

  @override
  Widget build(BuildContext context) {
    if (_error != null) return Center(child: Column(mainAxisSize: MainAxisSize.min, children: [Text(_error!, style: const TextStyle(color: Colors.red)), TextButton(onPressed: _loadList, child: const Text('Повторить'))]));

    if (_selectedNoteId != null) {
      if (_noteDetail == null) {
        _loadNote(_selectedNoteId!);
        return const Center(child: CircularProgressIndicator());
      }
      return RefreshIndicator(
        onRefresh: () => _loadNote(_selectedNoteId!),
        child: ListView(
          padding: const EdgeInsets.all(16),
          children: [
            TextButton.icon(icon: const Icon(Icons.arrow_back), label: const Text('Назад к заметкам'), onPressed: _goBack),
            const SizedBox(height: 8),
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(_noteDetail!['title']?.toString() ?? 'Заметка', style: Theme.of(context).textTheme.titleLarge),
                    const SizedBox(height: 8),
                    Text(_noteDetail!['body']?.toString() ?? ''),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 8),
            Text('Дополнения (${_comments.length})', style: Theme.of(context).textTheme.titleSmall),
            ..._comments.map<Widget>((c) => Card(
              margin: const EdgeInsets.only(top: 4),
              child: ListTile(
                title: Text(c['body']?.toString() ?? ''),
                subtitle: Text(c['created_at']?.toString() ?? ''),
              ),
            )),
            const SizedBox(height: 16),
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text('Добавить дополнение', style: Theme.of(context).textTheme.titleSmall),
                    TextField(
                      controller: _commentController,
                      decoration: const InputDecoration(hintText: 'Текст дополнения'),
                      maxLines: 2,
                    ),
                    const SizedBox(height: 8),
                    FilledButton(onPressed: _addComment, child: const Text('Добавить')),
                  ],
                ),
              ),
            ),
          ],
        ),
      );
    }

    return RefreshIndicator(
      onRefresh: _loadList,
      child: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          FilledButton.icon(icon: const Icon(Icons.add), label: const Text('Новая заметка'), onPressed: () => setState(() => _showForm = true)),
          if (_showForm) ...[
            const SizedBox(height: 16),
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    TextField(controller: _titleController, decoration: const InputDecoration(labelText: 'Заголовок')),
                    TextField(controller: _bodyController, decoration: const InputDecoration(labelText: 'Текст'), maxLines: 3),
                    const SizedBox(height: 8),
                    Row(
                      children: [
                        FilledButton(onPressed: _createNote, child: const Text('Создать')),
                        TextButton(onPressed: () => setState(() => _showForm = false), child: const Text('Отмена')),
                      ],
                    ),
                  ],
                ),
              ),
            ),
          ],
          const SizedBox(height: 16),
          ..._list.map<Widget>((n) {
            final map = n as Map<String, dynamic>;
            final id = map['id'];
            final title = map['title']?.toString() ?? '(без заголовка)';
            final body = (map['body']?.toString() ?? '').length > 80 ? '${(map['body'] as String).substring(0, 80)}...' : (map['body']?.toString() ?? '');
            return Card(
              margin: const EdgeInsets.only(bottom: 8),
              child: ListTile(
                title: Text(title),
                subtitle: Text(body),
                onTap: () => setState(() => _selectedNoteId = id),
              ),
            );
          }),
        ],
      ),
    );
  }
}
