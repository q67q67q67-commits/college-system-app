import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/auth_provider.dart';
import '../services/api_service.dart';

class DirectorScreen extends StatefulWidget {
  const DirectorScreen({super.key});

  @override
  State<DirectorScreen> createState() => _DirectorScreenState();
}

class _DirectorScreenState extends State<DirectorScreen> {
  Map<String, dynamic>? _data;
  List<dynamic> _posts = [];
  int? _selectedPostId;
  Map<String, dynamic>? _postDetail;
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

  Future<void> _load() async {
    try {
      final results = await Future.wait([
        ApiService.get('/api/director'),
        ApiService.get('/api/director/posts'),
      ]);
      final d = results[0];
      final ps = results[1];
      setState(() {
        _data = d != null ? Map<String, dynamic>.from(d as Map) : null;
        _posts = ps is List ? List<dynamic>.from(ps as List) : [];
        _error = null;
      });
    } catch (e) {
      setState(() {
        _data = null;
        _posts = [];
        _error = e.toString();
      });
    }
  }

  Future<void> _loadPost(int id) async {
    try {
      final res = await ApiService.get('/api/director/posts/$id');
      setState(() {
        _postDetail = res['post'] is Map ? Map<String, dynamic>.from(res['post'] as Map) : null;
        _comments = res['comments'] is List ? List<dynamic>.from(res['comments'] as List) : [];
        _error = null;
      });
    } catch (e) {
      setState(() {
        _postDetail = null;
        _comments = [];
        _error = e.toString();
      });
    }
  }

  Future<void> _createPost() async {
    final title = _titleController.text.trim();
    final body = _bodyController.text.trim();
    if (body.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Введите текст')));
      return;
    }
    try {
      await ApiService.post('/api/director/posts', {'title': title, 'body': body});
      _titleController.clear();
      _bodyController.clear();
      setState(() => _showForm = false);
      _load();
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(e.toString())));
    }
  }

  Future<void> _addComment() async {
    final body = _commentController.text.trim();
    if (body.isEmpty || _selectedPostId == null) return;
    try {
      await ApiService.post('/api/director/posts/$_selectedPostId/comments', {'body': body});
      _commentController.clear();
      _loadPost(_selectedPostId!);
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(e.toString())));
    }
  }

  void _goBack() {
    setState(() {
      _selectedPostId = null;
      _postDetail = null;
      _comments = [];
    });
    _load();
  }

  @override
  void initState() {
    super.initState();
    _load();
  }

  @override
  Widget build(BuildContext context) {
    final auth = context.watch<AuthProvider>();
    final isDirector = auth.role == 'director';

    if (_error != null) return Center(child: Column(mainAxisSize: MainAxisSize.min, children: [Text(_error!, style: const TextStyle(color: Colors.red)), TextButton(onPressed: _load, child: const Text('Повторить'))]));

    if (_selectedPostId != null) {
      if (_postDetail == null) {
        _loadPost(_selectedPostId!);
        return const Center(child: CircularProgressIndicator());
      }
      return RefreshIndicator(
        onRefresh: () => _loadPost(_selectedPostId!),
        child: ListView(
          padding: const EdgeInsets.all(16),
          children: [
            TextButton.icon(icon: const Icon(Icons.arrow_back), label: const Text('Назад к постам'), onPressed: _goBack),
            const SizedBox(height: 8),
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(_postDetail!['title']?.toString() ?? 'Пост', style: Theme.of(context).textTheme.titleLarge),
                    const SizedBox(height: 8),
                    Text(_postDetail!['body']?.toString() ?? ''),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 8),
            Text('Комментарии (${_comments.length})', style: Theme.of(context).textTheme.titleSmall),
            ..._comments.map<Widget>((c) => Card(
              margin: const EdgeInsets.only(top: 4),
              child: ListTile(
                title: Text(c['body']?.toString() ?? ''),
                subtitle: Text('${c['author_name'] ?? ''} · ${c['created_at'] ?? ''}'),
              ),
            )),
            const SizedBox(height: 16),
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text('Комментировать', style: Theme.of(context).textTheme.titleSmall),
                    TextField(
                      controller: _commentController,
                      decoration: const InputDecoration(hintText: 'Текст комментария'),
                      maxLines: 2,
                    ),
                    const SizedBox(height: 8),
                    FilledButton(onPressed: _addComment, child: const Text('Отправить')),
                  ],
                ),
              ),
            ),
          ],
        ),
      );
    }

    if (_data == null) return const Center(child: CircularProgressIndicator());
    final d = _data!;
    final phone = d['phone']?.toString() ?? '';

    return RefreshIndicator(
      onRefresh: _load,
      child: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(d['full_name']?.toString() ?? '', style: Theme.of(context).textTheme.titleLarge),
                  if (d['email']?.toString().isNotEmpty == true) Text(d['email'].toString()),
                  if (phone.isNotEmpty) FilledButton.icon(icon: const Icon(Icons.phone), label: Text('Написать: $phone'), onPressed: () {}),
                ],
              ),
            ),
          ),
          if (isDirector) ...[
            const SizedBox(height: 8),
            FilledButton.icon(icon: const Icon(Icons.add), label: const Text('Новый пост'), onPressed: () => setState(() => _showForm = true)),
            if (_showForm) ...[
              const SizedBox(height: 8),
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
                          FilledButton(onPressed: _createPost, child: const Text('Опубликовать')),
                          TextButton(onPressed: () => setState(() => _showForm = false), child: const Text('Отмена')),
                        ],
                      ),
                    ],
                  ),
                ),
              ),
            ],
          ],
          const SizedBox(height: 16),
          Text('Посты', style: Theme.of(context).textTheme.titleMedium),
          const SizedBox(height: 8),
          ...(_posts.isEmpty
              ? [const Card(child: Padding(padding: EdgeInsets.all(16), child: Text('Нет постов')))]
              : _posts.map<Widget>((p) {
                  final map = p as Map<String, dynamic>;
                  final id = map['id'];
                  final title = map['title']?.toString() ?? '(без заголовка)';
                  final body = (map['body']?.toString() ?? '').length > 80 ? '${(map['body'] as String).substring(0, 80)}...' : (map['body']?.toString() ?? '');
                  return Card(
                    margin: const EdgeInsets.only(bottom: 8),
                    child: ListTile(
                      title: Text(title),
                      subtitle: Text(body),
                      onTap: () => setState(() => _selectedPostId = id),
                    ),
                  );
                }).toList()),
        ],
      ),
    );
  }
}
