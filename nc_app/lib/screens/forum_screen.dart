import 'package:flutter/material.dart';
import '../services/api_service.dart';

class ForumScreen extends StatefulWidget {
  const ForumScreen({super.key});

  @override
  State<ForumScreen> createState() => _ForumScreenState();
}

String _truncate(String? s, int len) {
  final t = s ?? '';
  return t.length > len ? '${t.substring(0, len)}...' : t;
}

class _ForumScreenState extends State<ForumScreen> {
  List<dynamic> _list = [];
  String? _error;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final res = await ApiService.get('/api/forum/posts?limit=30');
      setState(() {
        _list = res is List ? List<dynamic>.from(res as List) : [];
        _list = _list.where((p) => p['parent_id'] == null).toList();
        _error = null;
      });
    } catch (e) {
      setState(() {
        _list = [];
        _error = e.toString();
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_error != null) return Center(child: Column(mainAxisSize: MainAxisSize.min, children: [Text(_error!, style: const TextStyle(color: Colors.red)), TextButton(onPressed: _load, child: const Text('Повторить'))]));
    return RefreshIndicator(
      onRefresh: _load,
      child: ListView.builder(
        padding: const EdgeInsets.all(16),
        itemCount: _list.length,
        itemBuilder: (_, i) {
          final p = _list[i] as Map<String, dynamic>;
          return Card(
            margin: const EdgeInsets.only(bottom: 8),
            child: ListTile(
              title: Text(p['title']?.toString() ?? '(без заголовка)'),
              subtitle: Text('${_truncate(p['body']?.toString(), 80)}\n${p['is_anonymous'] == true ? 'Анонимно' : p['author_name'] ?? ''}'),
              isThreeLine: true,
              onTap: () async {
                try {
                  final replies = await ApiService.get('/api/forum/posts?parent_id=${p['id']}&limit=50');
                  if (!context.mounted) return;
                  showModalBottomSheet(
                    context: context,
                    isScrollControlled: true,
                    builder: (_) => DraggableScrollableSheet(
                      initialChildSize: 0.7,
                      expand: false,
                      builder: (_, c) => ListView(
                        controller: c,
                        padding: const EdgeInsets.all(16),
                        children: [
                          Text(p['body']?.toString() ?? '', style: Theme.of(context).textTheme.bodyLarge),
                          const Divider(),
                          Text('Ответы (${replies is List ? replies.length : 0})', style: Theme.of(context).textTheme.titleSmall),
                          ...(replies is List ? (replies as List).map<Widget>((r) => ListTile(title: Text(r['body']?.toString() ?? ''), subtitle: Text(r['is_anonymous'] == true ? 'Анонимно' : r['author_name']?.toString() ?? ''))) : <Widget>[]),
                        ],
                      ),
                    ),
                  );
                } catch (_) {}
              },
            ),
          );
        },
      ),
    );
  }
}
