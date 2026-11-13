const sequence = (dir, count) =>
  Array.from({ length: count }, (_, i) => `${dir}/${i}.png`);
const digits = sequence('time', 10);
const small = sequence('info', 10);
const weekdays = sequence('week', 7);

WatchFace({
  build() {
    const ui = hmUI;
    const widget = ui.widget;
    const normal = ui.show_level.ONLY_NORMAL;
    const left = ui.align.LEFT;
    const image = (src, x, y) => ui.createWidget(widget.IMG, {
      src,
      x,
      y,
      show_level: normal
    });

    ui.createWidget(widget.FILL_RECT, {
      x: 0,
      y: 0,
      w: 390,
      h: 450,
      color: 0,
      show_level: normal
    });

    image('dial/left.png', 12, 102);
    image('dial/top.png', 72, 42);
    image('dial/bottom.png', 72, 354);
    image('dial/right.png', 324, 102);

    // 後に作ったウィジェットが前面に描かれる。文字盤内の文字を針で隠さないよう、針を先に作る。
    ui.createWidget(widget.TIME_POINTER, {
      hour_centerX: 195,
      hour_centerY: 225,
      hour_posX: 6,
      hour_posY: 102,
      hour_path: 'hand/hour.png',
      minute_centerX: 195,
      minute_centerY: 225,
      minute_posX: 8,
      minute_posY: 166,
      minute_path: 'hand/minute.png',
      second_centerX: 195,
      second_centerY: 225,
      second_posX: 6,
      second_posY: 178,
      second_path: 'hand/second.png',
      show_level: normal
    });

    ui.createWidget(widget.IMG_TIME, {
      hour_startX: 131,
      hour_startY: 113,
      hour_array: digits,
      hour_zero: 1,
      hour_space: 0,
      hour_follow: 0,
      hour_align: left,
      hour_unit_en: 'time/colon.png',
      hour_unit_sc: 'time/colon.png',
      hour_unit_tc: 'time/colon.png',
      minute_startX: 204,
      minute_startY: 113,
      minute_array: digits,
      minute_zero: 1,
      minute_space: 0,
      minute_follow: 0,
      minute_align: left,
      show_level: normal
    });

    ui.createWidget(widget.IMG_DATE, {
      month_startX: 149,
      month_startY: 288,
      month_zero: 1,
      month_space: 0,
      month_align: left,
      month_follow: 0,
      month_unit_en: 'info/slash.png',
      month_unit_sc: 'info/slash.png',
      month_unit_tc: 'info/slash.png',
      month_sc_array: small,
      month_tc_array: small,
      month_en_array: small,
      day_startX: 187,
      day_startY: 288,
      day_zero: 1,
      day_space: 0,
      day_align: left,
      day_follow: 0,
      day_sc_array: small,
      day_tc_array: small,
      day_en_array: small,
      show_level: normal
    });

    ui.createWidget(widget.IMG_WEEK, {
      x: 220,
      y: 286,
      week_en: weekdays,
      week_sc: weekdays,
      week_tc: weekdays,
      show_level: normal
    });

    ui.createWidget(widget.TEXT_IMG, {
      x: 145,
      y: 316,
      w: 100,
      h: 22,
      type: ui.data_type.BATTERY,
      font_array: small,
      h_space: 0,
      align_h: ui.align.CENTER_H,
      unit_en: 'info/percent.png',
      unit_sc: 'info/percent.png',
      unit_tc: 'info/percent.png',
      show_level: normal
    });
  }
});
