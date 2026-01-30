import 'package:flutter/material.dart';
import 'package:flutter/foundation.dart' show kIsWeb;
import 'package:file_picker/file_picker.dart';
import '../utils/download_helper.dart';
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
  bool _uploading = false;

  Future<void> _pickAndUpload() async {
    final result = await FilePicker.platform.pickFiles(withData: true);
    if (result == null || result.files.isEmpty) return;
    final file = result.files.first;
    final bytes = file.bytes;
    if (bytes == null) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Не удалось прочитать файл')));
      return;
    }
    setState(() => _uploading = true);
    try {
      await ApiService.uploadFile('/api/files/upload', bytes, file.name);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Файл загружен')));
        _load();
      }
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(e.toString())));
    } finally {
      if (mounted) setState(() => _uploading = false);
    }
  }

  Future<void> _download(int id, String filename) async {
    try {
      final r = await ApiService.getBytes('/api/files/$id/download');
      if (kIsWeb) {
        if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Скачано ${r.bodyBytes.length} байт. Для скачивания откройте в браузере /app/.')));
      } else {
        await saveAndOpenFile(r.bodyBytes, filename);
        if (mounted) ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Файл открыт')));
      }
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(e.toString())));
    }
  }

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final res = await ApiService.get('/api/files');
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
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text('Использовано: ${(_used / 1024 / 1024).toStringAsFixed(1)} МБ из ${(_limit / 1024 / 1024 / 1024).toStringAsFixed(1)} ГБ'),
                  const SizedBox(height: 8),
                  FilledButton.icon(
                    onPressed: _uploading ? null : _pickAndUpload,
                    icon: _uploading ? const SizedBox(width: 16, height: 16, child: CircularProgressIndicator(strokeWidth: 2)) : const Icon(Icons.upload),
                    label: Text(_uploading ? 'Загрузка...' : 'Загрузить файл'),
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 8),
          ..._files.map<Widget>((f) {
            final map = f as Map<String, dynamic>;
            final id = map['id'];
            final filename = map['filename']?.toString() ?? map['path']?.toString() ?? 'file';
            return Card(
              margin: const EdgeInsets.only(bottom: 4),
              child: ListTile(
                title: Text(filename),
                subtitle: Text('${((map['size_bytes'] as num?)?.toInt() ?? 0) / 1024} КБ'),
                trailing: Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    IconButton(
                      icon: const Icon(Icons.download),
                      onPressed: () => _download(id, filename),
                    ),
                    IconButton(
                      icon: const Icon(Icons.delete),
                      onPressed: () async {
                        try {
                          await ApiService.delete('/api/files/$id');
                          if (context.mounted) _load();
                        } catch (_) {}
                      },
                    ),
                  ],
                ),
              ),
            );
          }),
        ],
      ),
    );
  }
}
