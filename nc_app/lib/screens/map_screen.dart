import 'package:flutter/material.dart';

class MapScreen extends StatelessWidget {
  const MapScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        Card(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('Карта здания', style: Theme.of(context).textTheme.titleLarge),
                const SizedBox(height: 12),
                const Text('Схема этажей и навигация по корпусу. Адрес: ул. Жандосова, 55, г. Алматы, Казахстан, 050035.'),
                const SizedBox(height: 24),
                Container(
                  height: 200,
                  decoration: BoxDecoration(color: Colors.grey[300], borderRadius: BorderRadius.circular(8)),
                  child: const Center(child: Text('Карта и план эвакуации\n(размещаются администрацией)', textAlign: TextAlign.center)),
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }
}
