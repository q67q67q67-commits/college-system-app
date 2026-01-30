import 'package:flutter/material.dart';
import '../services/api_service.dart';

class LibraryScreen extends StatefulWidget {
  const LibraryScreen({super.key});

  @override
  State<LibraryScreen> createState() => _LibraryScreenState();
}

class _LibraryScreenState extends State<LibraryScreen> {
  List<dynamic> _list = [];
  String? _error;
  final _query = TextEditingController();

  @override
  void initState() {
    super.initState();
    _load();
  }

  @override
  void dispose() {
    _query.dispose();
    super.dispose();
  }

  Future<void> _load() async {
    try {
      final q = _query.text.trim();
      final path = q.isEmpty ? '/api/library/books?limit=20' : '/api/library/books?q=${Uri.encodeComponent(q)}&limit=20';
      final res = await ApiService.get(path);
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

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.all(8),
          child: Row(
            children: [
              Expanded(child: TextField(controller: _query, decoration: const InputDecoration(labelText: 'Поиск', border: OutlineInputBorder()), onSubmitted: (_) => _load())),
              IconButton(icon: const Icon(Icons.search), onPressed: _load),
            ],
          ),
        ),
        Expanded(
          child: _error != null
              ? Center(child: Column(mainAxisSize: MainAxisSize.min, children: [Text(_error!, style: const TextStyle(color: Colors.red)), TextButton(onPressed: _load, child: const Text('Повторить'))]))
              : RefreshIndicator(
                  onRefresh: _load,
                  child: ListView.builder(
                    padding: const EdgeInsets.all(16),
                    itemCount: _list.length,
                    itemBuilder: (_, i) {
                      final b = _list[i] as Map<String, dynamic>;
                      return Card(
                        margin: const EdgeInsets.only(bottom: 8),
                        child: ListTile(
                          title: Text(b['title']?.toString() ?? ''),
                          subtitle: Text('${b['author'] ?? ''} · Доступно: ${b['available'] ?? 0}'),
                          trailing: FilledButton(
                            onPressed: () async {
                              try {
                                await ApiService.post('/api/library/reservations', {'book_id': b['id']});
                                if (context.mounted) ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Бронь создана')));
                              } catch (e) {
                                if (context.mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(e.toString())));
                              }
                            },
                            child: const Text('Забронировать'),
                          ),
                        ),
                      );
                    },
                  ),
                ),
        ),
      ],
    );
  }
}
