import 'package:flutter/material.dart';
import '../services/api_service.dart';

class GradesScreen extends StatefulWidget {
  const GradesScreen({super.key});

  @override
  State<GradesScreen> createState() => _GradesScreenState();
}

class _GradesScreenState extends State<GradesScreen> {
  List<dynamic> _list = [];
  double? _gpa;
  String? _error;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final grades = await ApiService.get('/api/grades');
      final gpaRes = await ApiService.get('/api/grades/gpa');
      setState(() {
        _list = grades is List ? List<dynamic>.from(grades as List) : [];
        _gpa = (gpaRes['gpa'] as num?)?.toDouble();
        _error = null;
      });
    } catch (e) {
      setState(() {
        _list = [];
        _gpa = null;
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
          if (_gpa != null) Card(child: Padding(padding: const EdgeInsets.all(16), child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [Text('Средний балл (GPA): $_gpa', style: Theme.of(context).textTheme.titleLarge), const SizedBox(height: 4), Text('Оценки по 100-балльной шкале', style: Theme.of(context).textTheme.bodySmall)]))),
          const SizedBox(height: 8),
          ..._list.map<Widget>((g) => Card(
                margin: const EdgeInsets.only(bottom: 4),
                child: ListTile(
                  title: Text(g['subject']?.toString() ?? ''),
                  subtitle: Text('${g['grade_date']} — ${g['grade']} баллов${g['comment'] != null ? ' · ${g['comment']}' : ''}'),
                ),
              )),
        ],
      ),
    );
  }
}
