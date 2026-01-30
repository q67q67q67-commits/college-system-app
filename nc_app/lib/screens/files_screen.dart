import 'package:flutter/material.dart';
import '../services/api_service.dart';

class FilesScreen extends StatefulWidget {
  const FilesScreen({super.key});

  @override
  State<FilesScreen> createState() => _FilesScreenState();
}

class _FilesScreenState extends State<FilesScreen> {
  List<dynamic> _files = [];
  int _used = 0;
  int _limit = 0;
  String? _error;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final res = await ApiService.get('/api/files') as Map<String, dynamic>;
      setState(() {
        _files = res['files'] is List ? res['files'] as List : [];
        _used = (res['storage_used'] as num?)?.toInt() ?? 0;
        _limit = (res['storage_limit'] as num?)?.toInt() ?? 0;
        _error = null;
      });
    } catch (e) {
      setState(() {
        _files = [];
        _error = e.toString();
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_error != null) return Center(child: Column(mainAxisSize: MainAxisSize.min, children: [Text(_error!, style: const TextStyle(color: Colors.red)), TextButton(onPressed: _load, child: const Text('Повторить'))]));
    return RefreshIndicator(
      onRefresh: _load,
      child: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Text('Использовано: ${(_used / 1024 / 1024).toStringAsFixed(1)} МБ из ${(_limit / 1024 / 1024 / 1024).toStringAsFixed(1)} ГБ'),
            ),
          ),
          const SizedBox(height: 8),
          ..._files.map<Widget>((f) {
            final map = f as Map<String, dynamic>;
            return Card(
              margin: const EdgeInsets.only(bottom: 4),
              child: ListTile(
                title: Text(map['filename']?.toString() ?? map['path']?.toString() ?? ''),
                subtitle: Text('${((map['size_bytes'] as num?)?.toInt() ?? 0) / 1024} КБ'),
                trailing: IconButton(
                  icon: const Icon(Icons.delete),
                  onPressed: () async {
                    try {
                      await ApiService.delete('/api/files/${map['id']}');
                      if (context.mounted) _load();
                    } catch (_) {}
                  },
                ),
              ),
            );
          }),
        ],
      ),
    );
  }
}
