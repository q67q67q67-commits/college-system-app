import 'package:flutter/material.dart';
import '../services/api_service.dart';

class DirectorScreen extends StatefulWidget {
  const DirectorScreen({super.key});

  @override
  State<DirectorScreen> createState() => _DirectorScreenState();
}

class _DirectorScreenState extends State<DirectorScreen> {
  Map<String, dynamic>? _data;
  String? _error;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final res = await ApiService.get('/api/director') as Map<String, dynamic>;
      setState(() {
        _data = res;
        _error = null;
      });
    } catch (e) {
      setState(() {
        _data = null;
        _error = e.toString();
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_error != null) return Center(child: Column(mainAxisSize: MainAxisSize.min, children: [Text(_error!, style: const TextStyle(color: Colors.red)), TextButton(onPressed: _load, child: const Text('Повторить'))]));
    if (_data == null) return const Center(child: CircularProgressIndicator());
    final d = _data!;
    final phone = d['phone']?.toString() ?? '';
    return Center(
      child: Card(
        margin: const EdgeInsets.all(24),
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              const Text('Профиль директора', style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
              const SizedBox(height: 16),
              Text(d['full_name']?.toString() ?? '', style: Theme.of(context).textTheme.titleLarge),
              if (d['email']?.toString().isNotEmpty == true) Text(d['email'].toString()),
              if (phone.isNotEmpty) ...[
                const SizedBox(height: 8),
                FilledButton.icon(
                  icon: const Icon(Icons.phone),
                  label: Text('Написать директору: $phone'),
                  onPressed: () {},
                ),
              ],
            ],
          ),
        ),
      ),
    );
  }
}
