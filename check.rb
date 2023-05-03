fname = ARGV[0]
incidence = Hash.new { |h, k| h[k] = {} }
File.read(fname).split("\n").select { |x| x[0] == 'e' }.each do |line|
	_, v1, v2 = line.split
	incidence[v1][v2] = 1
	incidence[v2][v1] = 1
end

res = incidence.keys.select do |k|
	incidence[k].count != 4
end
puts incidence.count
puts res
