#!/usr/bin/env perl
#
# Script to build TOC for README.md using https://github.com/ekalinin/github-markdown-toc
#

use strict;
use warnings;

my $file = scalar @ARGV > 0 ? $ARGV[1] : 'README.md';
die("No such markdown file \"$file\"\n") unless -f $file;

my $num = 0;
my $toc = join("\n", grep {$_} map {
    $num++;
    chomp;
    $num > 1 ? substr($_, 4) : ""; # strip first line, reduce indent on all
} split(/\n/, `cat $file | gh-md-toc -`));

open my $in, '<', $file or die("Could not open \"$file\" for read: $!\n");
my $tmp = "$file.tmp__";
open my $out, '>', $tmp or die("Could not open \"$tmp\" for write: $!\n");
my $state = '';
my $written = 0;
my $closed = 0;
foreach my $line(<$in>) {
    chomp($line);
    if ($state eq 'toc') {
        if ($line =~ /<!--\s*TOC\s*END\s*-->/) {
            $closed++;
            $state = '';
            print $out "\n$line\n";
        }
    } elsif ($line =~ /<!--\s*TOC\s*START\s*-->/) {
        print $out "$line\n\n";
        print $out "$toc\n";
        $state = 'toc';
        $written++;
    } else {
        print $out "$line\n";
    }
}
close $in;
close $out;
#`mv $tmp $file`;

my @errs;
push @errs, "TOC was not written. Start marker \"<!-- TOC START -->\" was not found!" unless $written;
push @errs, "TOC was never closed. End marker \"<!-- TOC END -->\" was not found!" unless $closed;
die("Failed:\n * ". join("\n * ", @errs). "\n") if @errs;

`mv $tmp $file`;
print "TOC extracted and written to \"$file\"\n";
