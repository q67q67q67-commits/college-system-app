import 'package:flutter/material.dart';
import '../../services/api_client.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key, required this.apiToken, required this.onLogout});

  final String apiToken;
  final VoidCallback onLogout;

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  late final ApiClient _api;

  @override
  void initState() {
    super.initState();
    _api = ApiClient(ApiClient.defaultBaseUrl, token: widget.apiToken);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Главная'),
        actions: [
          IconButton(
            icon: const Icon(Icons.logout),
            onPressed: () async => widget.onLogout(),
          ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: () async => setState(() {}),
        child: ListView(
          padding: const EdgeInsets.all(16),
          children: [
            _EventsBlock(api: _api),
            const SizedBox(height: 24),
            _DirectorBlock(),
            const SizedBox(height: 24),
            _ScheduleBlock(api: _api),
          ],
        ),
      ),
    );
  }
}

class _EventsBlock extends StatelessWidget {
  const _EventsBlock({required this.api});

  final ApiClient api;

  @override
  Widget build(BuildContext context) {
    return FutureBuilder<List<dynamic>>(
      future: api.getList('/api/events?limit=10'),
      builder: (context, snapshot) {
        if (snapshot.connectionState == ConnectionState.waiting) {
          return const Card(
            child: Padding(
              padding: EdgeInsets.all(24),
              child: Center(child: CircularProgressIndicator()),
            ),
          );
        }
        if (snapshot.hasError) {
          return Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Text(
                'Новости: ${snapshot.error}',
                style: TextStyle(color: Colors.red.shade700),
              ),
            ),
          );
        }
        final list = snapshot.data ?? [];
        return Card(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('Новости колледжа', style: Theme.of(context).textTheme.headlineSmall),
                const SizedBox(height: 12),
                if (list.isEmpty)
                  const Text('Пока нет новостей.')
                else
                  ...list.take(5).map((e) {
                    final map = e as Map<String, dynamic>;
                    final title = map['title'] as String? ?? '';
                    final desc = map['description'] as String? ?? '';
                    return Padding(
                      padding: const EdgeInsets.only(bottom: 8),
                      child: ListTile(
                        contentPadding: EdgeInsets.zero,
                        title: Text(title),
                        subtitle: desc.isNotEmpty ? Text(desc, maxLines: 2, overflow: TextOverflow.ellipsis) : null,
                      ),
                    );
                  }),
              ],
            ),
          ),
        );
      },
    );
  }
}

class _DirectorBlock extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Профиль директора', style: Theme.of(context).textTheme.headlineSmall),
            const SizedBox(height: 12),
            const Row(
              children: [
                CircleAvatar(radius: 32, child: Icon(Icons.person, size: 32)),
                SizedBox(width: 16),
                Expanded(
                  child: Text('Краткое описание директора колледжа. Переход на полную страницу — в разработке.'),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

class _ScheduleBlock extends StatelessWidget {
  const _ScheduleBlock({required this.api});

  final ApiClient api;

  @override
  Widget build(BuildContext context) {
    return FutureBuilder<List<dynamic>>(
      future: api.getList('/api/schedule'),
      builder: (context, snapshot) {
        if (snapshot.connectionState == ConnectionState.waiting) {
          return const Card(
            child: Padding(
              padding: EdgeInsets.all(24),
              child: Center(child: CircularProgressIndicator()),
            ),
          );
        }
        if (snapshot.hasError) {
          return Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Text(
                'Расписание: ${snapshot.error}',
                style: TextStyle(color: Colors.red.shade700),
              ),
            ),
          );
        }
        final list = snapshot.data ?? [];
        final today = DateTime.now().weekday;
        final todayLessons = list.where((e) {
          final d = e['day_of_week'];
          final n = d is int ? d : (d as num?)?.toInt();
          return n == today;
        }).toList();
        todayLessons.sort((a, b) {
          final t1 = a['start_time'] as String? ?? '';
          final t2 = b['start_time'] as String? ?? '';
          return t1.compareTo(t2);
        });

        return Card(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('Сегодняшние занятия', style: Theme.of(context).textTheme.headlineSmall),
                const SizedBox(height: 12),
                if (todayLessons.isEmpty)
                  const Text('Сегодня выходной.')
                else
                  ...todayLessons.map((e) {
                    final map = e as Map<String, dynamic>;
                    final subject = map['subject'] as String? ?? '';
                    final room = map['room'] as String? ?? '';
                    final start = map['start_time'] as String? ?? '';
                    return ListTile(
                      contentPadding: EdgeInsets.zero,
                      title: Text(subject),
                      subtitle: Text('$room · $start'),
                    );
                  }),
              ],
            ),
          ),
        );
      },
    );
  }
}
